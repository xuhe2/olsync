package internal

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/goccy/go-yaml"
)

// Config 是整个配置文件的顶级结构体。
type Config struct {
	Overleaf OverleafConfig `yaml:"overleaf"`
}

// OverleafConfig 包含了 Overleaf 服务的配置信息。
type OverleafConfig struct {
	BaseURL string `yaml:"baseUrl"`
	Cookies []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"cookies"`
}

// ParseConfig 从指定路径解析 YAML 配置文件。
func ParseConfig(reader io.Reader) (*Config, error) {
	// 读取 reader 中的全部数据
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("error reading from reader: %v", err)
		return nil, err
	}

	// 初始化一个 Config 结构体
	cfg := &Config{}

	// 使用 goccy/go-yaml 库解析 YAML 数据
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Printf("error unmarshalling config: %v", err)
		return nil, err
	}

	return cfg, nil
}

// 从文件路径解析配置文件
func ParseConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	return ParseConfig(file)
}

func (c *Config) GetCookies() []*http.Cookie {
	var httpCookies []*http.Cookie
	for _, cookie := range c.Overleaf.Cookies {
		httpCookie := &http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
		httpCookies = append(httpCookies, httpCookie)
	}
	return httpCookies
}

func (c *Config) GetBaseURL() url.URL {
	url, err := url.Parse(c.Overleaf.BaseURL)
	if err != nil {
		log.Printf("error parsing base url: %v", err)
		log.Fatalf("invalid base url: %s", c.Overleaf.BaseURL)
	}
	return *url
}
