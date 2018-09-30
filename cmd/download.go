package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"bufio"
	"strings"
	"sync"
	"time"
	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"  
)

var (
	Artist string
	Title string
	Filename string
	Target string
	Concurrency int
	StringList string
)
  
func hash(s string) string {
	return fmt.Sprintf("%x",s) 
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func countLines(filename string) int {
	i := 0
	file, err := os.Open(filename)
	check(err)
	
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		i++
	}

	return i
}
  
func Download(title string, artist string){ 
  
	_ = os.Mkdir(artist,0777)
  
	filename := artist +"/"+ title+".mp3"
	filehash := hash(filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
  
		fileStaging := artist + "/"+ filehash
		if _, err := os.Stat(fileStaging); os.IsNotExist(err) {
			_ = os.Mkdir(fileStaging,0777)
		}

		cmd := exec.Command("song","-y",artist +" - "+title) 
		cmd.Dir = fileStaging
		err := cmd.Run()
		check(err)

		files, err := ioutil.ReadDir(fileStaging)
		check(err)

		for _, f := range files {
			err := os.Rename(fileStaging+"/"+f.Name(),filename)
			check(err)

			time.Sleep(1000)
			err = os.RemoveAll(fileStaging)
			check(err)
		}     
	} 
}

func DownloadFromString(str string){
	throttle := make(chan bool,Concurrency)
	defer close(throttle)

	var wg sync.WaitGroup

	lines := strings.Split(str,"\n")
	linecount := len(lines)

	bar := progressbar.New(linecount)
	bar.Add(1)

	for i := 0; i < len(lines); i++  {
		s := strings.Split(lines[i]," - ")

		wg.Add(1)
		go func(){
			throttle <- true
			Download(s[1],s[0])
			<-throttle
			bar.Add(1)
			wg.Done() 
		}()
	}

	wg.Wait()
	fmt.Println("")
}

func DownloadFromFile(filename string){
	throttle := make(chan bool,Concurrency)
	defer close(throttle)

	var wg sync.WaitGroup

	file, err := os.Open(filename)
	check(err)

	defer file.Close()

	linecount := countLines(filename)

	bar := progressbar.New(linecount+1)
	bar.Add(1)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line," - ")

		wg.Add(1)
		go func(){
			throttle <- true
			Download(s[1],s[0])
			<-throttle
			bar.Add(1)
			wg.Done() 
		}()
	}
	
	wg.Wait()
	fmt.Println("")
}

func init() {
	donwload.Flags().StringVarP(&Target, "target", "t", "", "Where do download songs")	

	single.Flags().StringVarP(&Artist, "artist", "a", "", "The artist of the song")
	single.Flags().StringVarP(&Title, "title", "t", "", "The title of the song")

	bulk.Flags().StringVarP(&Filename, "filename", "f", "", "The source file for downloading songs. Each line must be formatted like 'Artist - Title'")
	bulk.Flags().IntVarP(&Concurrency, "concurrency", "c", 1, "How many files to download concurrently")
	bulk.Flags().StringVarP(&StringList,"string-list","s","","The multiline variable containing the song list")

	donwload.AddCommand(single)
	donwload.AddCommand(bulk)
  	rootCmd.AddCommand(donwload)
}

var donwload = &cobra.Command{
  Use:   "download",
  Short: "download songs",
  Long:  `A command for performing song downloads`,
  Run: func(cmd *cobra.Command, args []string) {
    
  },
}

var bulk = &cobra.Command{
	Use:   "bulk",
	Short: "bulk song download",
	Long:  `A command for performing bulk song download with file`,
	Run: func(cmd *cobra.Command, args []string) {

		if (Target!=""){
			err := os.Chdir(Target)
			check(err)
		}

		if (StringList!=""){
			DownloadFromString(strings.Trim(StringList," "))
		} else {
			DownloadFromFile(Filename)
		}
	},
}

var single = &cobra.Command{
	Use:   "single",
	Short: "single song download",
	Long:  `A command for performing a single song download`,
	Run: func(cmd *cobra.Command, args []string) {

		if (Target!=""){
			err := os.Chdir(Target)
			check(err)
		}

		Download(Title, Artist)
	},
}