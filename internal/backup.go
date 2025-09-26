package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// 将非法字符替换为 _
func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_", "|", "_",
	)
	return replacer.Replace(name)
}

type BackupManager interface {
	BackupProject(project Project) error
	RunBackup(projects []Project) error
}

type ZipBackupManager struct {
	Client *OLClient
	Config BackupConfig
}

func NewZipBackupManager(client *OLClient, config BackupConfig) *ZipBackupManager {
	return &ZipBackupManager{
		Client: client,
		Config: config,
	}
}

// BackupProject 下载并备份单个项目
func (z *ZipBackupManager) BackupProject(project Project) error {
	reader, err := z.Client.DownloadProjectZip(project)
	if err != nil {
		return fmt.Errorf("下载项目 %s 失败: %w", project.Name, err)
	}
	defer reader.Close()

	// 生成当前备份的子目录，例如 ./Backup/2025-09-25_21-37-00
	timestampDir := time.Now().Format("2006-01-02_15-04-05")
	backupDir := filepath.Join(z.Config.Path, timestampDir)
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 文件名只用项目名
	safeName := sanitizeFileName(project.Name)
	fullPath := filepath.Join(backupDir, safeName+".zip")

	log.Printf("备份项目 %s 到 %s", project.Name, fullPath)

	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		// 删除已创建的备份文件
		os.Remove(fullPath)
		log.Printf("因为写入备份文件失败，已删除: %s", fullPath)

		return fmt.Errorf("写入备份文件失败: %w", err)
	}

	return nil
}

// RunBackup 批量备份
func (z *ZipBackupManager) RunBackup(projects []Project) error {
	schedule := strings.TrimSpace(z.Config.Schedule)
	if schedule == "" {
		log.Println("未设置定时任务，执行一次性备份")
		return z.runOnce(projects)
	}

	c := cron.New(cron.WithSeconds()) // 允许解析到秒（更精细）
	_, err := c.AddFunc(schedule, func() {
		log.Printf("触发定时备份任务: %s", schedule)
		if err := z.runOnce(projects); err != nil {
			log.Printf("定时备份失败: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("解析 Cron 表达式失败: %w", err)
	}

	log.Printf("启动定时备份任务, Cron 表达式: %s", schedule)
	c.Start()

	// 阻塞主线程，保持程序常驻
	select {}
}

func (z *ZipBackupManager) runOnce(projects []Project) error {
	var totalErr error
	for _, p := range projects {
		if err := z.BackupProject(p); err != nil {
			log.Printf("备份项目 %s 失败: %v", p.Name, err)
			totalErr = errors.Wrap(totalErr, err.Error()) // 记录所有错误，最后统一返回, 防止部分失败导致全部失败(容错)
		}
	}
	if totalErr != nil { // 如果存在错误，则返回
		return totalErr
	}

	// 清理旧备份文件夹
	// 只有全部成功写入文件后，才进行清理
	z.cleanupOldBackups()

	return nil
}

// 清理旧备份（按时间目录排序）
func (z *ZipBackupManager) cleanupOldBackups() {
	if z.Config.KeepLast <= 0 {
		return
	}

	dirs, err := os.ReadDir(z.Config.Path)
	if err != nil {
		log.Printf("读取备份目录失败: %v", err)
		return
	}

	// 仅保留文件夹
	var backupDirs []string
	for _, d := range dirs {
		if d.IsDir() {
			backupDirs = append(backupDirs, filepath.Join(z.Config.Path, d.Name()))
		}
	}

	if len(backupDirs) <= z.Config.KeepLast {
		return
	}

	sort.Slice(backupDirs, func(i, j int) bool {
		fi, _ := os.Stat(backupDirs[i])
		fj, _ := os.Stat(backupDirs[j])
		return fi.ModTime().Before(fj.ModTime())
	})

	toDelete := backupDirs[:len(backupDirs)-z.Config.KeepLast]
	for _, d := range toDelete {
		log.Printf("删除旧备份目录: %s", d)
		os.RemoveAll(d)
	}
}
