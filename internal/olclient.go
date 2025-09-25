package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	HomePageUrlPattern           = "/project/"
	DownloadProjectZipUrlPattern = "/project/%s/download/zip"
)

// OLClient 封装了与 Overleaf API 通信的客户端。
// 它维护一个 HTTP 客户端和一个会话 cookie。
type OLClient struct {
	client *http.Client

	projectPageUrl url.URL
	cookies        []*http.Cookie
}

// NewOLClient 创建并返回一个新的 OLClient 实例。
func NewOLClient() *OLClient {
	client := &OLClient{
		client: &http.Client{
			Timeout: 30 * time.Second, // 设置超时，防止请求挂起
		},
		projectPageUrl: url.URL{
			Scheme: "https",
			Host:   "www.overleaf.com",
		},
		cookies: make([]*http.Cookie, 0),
	}
	return client
}

func (c *OLClient) WithProjectPageUrl(url url.URL) *OLClient {
	c.projectPageUrl = url
	return c
}

// 通过 AOP 方式设置 Cookie（不会直接写入 Jar，而是每次请求动态注入）
func (c *OLClient) WithCookies(cookies []*http.Cookie) *OLClient {
	c.cookies = cookies
	// 包装 Transport，加一层拦截器
	baseTransport := c.client.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	c.client.Transport = &cookieInjectorTransport{
		base:    baseTransport,
		cookies: cookies,
	}
	return c
}

func (c *OLClient) WithTransport(transport http.RoundTripper) *OLClient {
	c.client.Transport = transport
	return c
}

// 需要手动Close io.ReadCloser
func (c *OLClient) VisitHomePage() (io.ReadCloser, error) {
	homePageUrl, err := url.JoinPath(c.projectPageUrl.String(), HomePageUrlPattern)
	if err != nil {
		return nil, err
	}
	log.Printf("Visit Home Page: %s", homePageUrl)
	resp, err := c.client.Get(homePageUrl)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *OLClient) GetProjects() []Project {
	reader, err := c.VisitHomePage()
	if err != nil {
		return nil
	}
	defer reader.Close()

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// 使用 CSS 选择器找到 `<meta name="ol-prefetchedProjectsBlob">` 元素。
	// 选择器 `meta[name="..."]` 是定位特定 meta 标签的标准方法。
	selection := doc.Find(`meta[name="ol-prefetchedProjectsBlob"]`)

	// 检查是否找到了该元素。
	if selection.Length() == 0 {
		log.Println("未找到 meta 标签 'ol-prefetchedProjectsBlob'")
		fmt.Println(doc.Html())
		return []Project{}
	}

	// 从找到的元素中获取 `content` 属性的值。
	content, exists := selection.Attr("content")
	if !exists {
		log.Println("meta 标签 'ol-prefetchedProjectsBlob' 没有 content 属性")
		return []Project{}
	}

	projectInfos, err := ParsePrefetchedProjectsBlob(content)
	if err != nil {
		return []Project{}
	}

	return projectInfos.Projects
}

func (c *OLClient) DownloadProjectZip(project Project) (io.ReadCloser, error) {
	downloadUrl, err := url.JoinPath(c.projectPageUrl.String(), fmt.Sprintf(DownloadProjectZipUrlPattern, project.ID))
	if err != nil {
		return nil, err
	}

	log.Printf("Download File %s zip from %s", project.Name, downloadUrl)

	resp, err := c.client.Get(downloadUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("下载项目失败，状态码：%d", resp.StatusCode)
	}

	return resp.Body, nil
}
