package main

import (
  "fmt"
  "flag"
  "os"
  "os/signal"
  "net/http"
  "io/ioutil"
  "path/filepath"
  "time"
  "strings"
  urllib "net/url"
  "github.com/robfig/cron"
)

func exists(path string) (bool) {
  _, err := os.Stat(path)
  if err == nil { 
    return true
  }
  return false
}

func scrapeUrl(url string) ([]byte, error) {
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return body, err
}

func scrape(workspace string, urls []string, layout string) {
  now := time.Now()
  folder := now.UTC().Format(layout)
  workdir := filepath.Clean(filepath.Join(workspace, folder))
  if !exists(workdir) {
    if err := os.MkdirAll(workdir, 0777); err != nil {
      panic(err)
    }
  }
  
  for _, url := range urls {
    content, err := scrapeUrl(url)
    if err != nil {
      fmt.Printf("Unable to scrape %s - %v\n", url, err)
      continue
    }
    
    u, err := urllib.Parse(url)
    if err != nil {
      panic(err)
    }    
    filePath := filepath.Join(workdir, u.Host, u.Path)
    basePath := filepath.Dir(filePath)
    if !exists(basePath) {
      if err := os.MkdirAll(basePath, 0777); err != nil {
        panic(err)
      }
    }
    if err = ioutil.WriteFile(filePath, content, 0644); err != nil { 
      panic(err) 
    }
    fmt.Printf("%v: fetched %s => %d byte(s) written into %s\n", now, url, len(content), filePath)
  }
}

func createWorkspaceIfNeeded(workspacePath string) {
  if !exists(workspacePath) {
    if err := os.MkdirAll(workspacePath, 0777); err != nil {
      panic(err)
    }
  }
}

func main() {
  versionPtr := flag.Bool("version", false, "prints version")
  testPtr := flag.Bool("test", false, "test scraping")
  workspacePtr := flag.String("workspace", "./workspace", "workspace directory")
  urlsPtr := flag.String("urls", "", "urls to scrape (coma delimited)")
  intervalPtr := flag.String("interval", "@hourly", "fetch interval, see github.com/robfig/cron")
  layoutPtr := flag.String("layout", "2006-01-02_15_04_05", "folder filename layout")
  
  flag.Parse()
  
  if *versionPtr {
    fmt.Println("scrap 0.1")
    return
  }
  
  if len(*urlsPtr)==0 {
    fmt.Printf("specify at least one url via --urls param\n")
    os.Exit(4)
  }
  urlsToScrape := strings.Split(*urlsPtr, ",")
  
  dir, _ := os.Getwd()
  workspace := *workspacePtr 
  if !filepath.IsAbs(*workspacePtr) {
    workspace = filepath.Clean(filepath.Join(dir, *workspacePtr))
  }
  createWorkspaceIfNeeded(workspace)
  
  if *testPtr {
    // test scrapping
    scrape(workspace, urlsToScrape, *layoutPtr)
  } else {
    // use cron and wait for signal to exit
    c := cron.New()
    c.AddFunc(*intervalPtr, func() { scrape(workspace, urlsToScrape, *layoutPtr) })
    c.Start()
  
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt, os.Kill)
    signal := <-signals
    fmt.Println("Got signal:", signal)
  }
}