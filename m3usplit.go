//==================== m3usplit =============================================
// The purpose of this program is to take an m3u file (playlist of audio files)
// and split it into multiple m3u playlists according to a size.
//
// The main use case would be to take a use playlist and split it into 700M
// m3u lists that would fit on a standard CD-ROM.
// 
// The new m3u files are placed in the same directory as the original m3u file
// and are named <original_name>-1.m3u, <original_name>-2.m3u, etc.
// 
//
//============================================================================


package main
import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"        
    "github.com/dj80hd/dj80hdutil"
)

func exitWithUsage() {
        programName := os.Args[0]
        //FIXME - Better way to do this like inline docstrings ?
        fmt.Printf("Splits an m3u file into many others given a length in megabytes\n")
        fmt.Printf("USAGE %q filename size_in_megabytes\n",programName)
        fmt.Printf("\n")                                         
        fmt.Printf("EXAMPLE: to split an m3u file into 700M segments:\n")
        fmt.Printf("%q foo/bar/baz.m3u 700",programName)
        os.Exit(1)                         
}

func main() {
    args := os.Args[1:] //ignore program name
    if len(args) != 2 {

        exitWithUsage()
    }

    fileCount := 0
    m3uFilename := args[0]
    m3uSplitLength,_ := strconv.ParseInt(args[1],0,64)
    m3uSplitLength = m3uSplitLength * 1024 * 1024 //700 becomes 700M

    //Data Structure to hold m3u files that are split
    var m3uSplitFileContents []string
    m3uSplitFileContents = append(m3uSplitFileContents,"")
    m3uSplitFileContentsIndex := 0

    //FIXME - Check valid m3u file and valid split length

    m3uFilenameDir := filepath.Dir(m3uFilename)

    fmt.Printf("YOU TYPED %q and %i\n",m3uFilename,m3uSplitLength)

    mp3Files := dj80hdutil.FileToArray(m3uFilename)

    currentSize := int64(0)        

    for i := 0; i < len(mp3Files); i++ {
        v := mp3Files[i]
        //resolve the name to the original m3u file
        resolvedFilename := filepath.Join(m3uFilenameDir,v) 

        fi,err := os.Stat(resolvedFilename)
        if err != nil {
            fmt.Printf("WARNING: Could not stat file %q\n",v)
            continue
        }

        size := fi.Size()
        if (currentSize + size) > m3uSplitLength {
            m3uSplitFileContentsIndex = m3uSplitFileContentsIndex + 1
            m3uSplitFileContents = append(m3uSplitFileContents,"")
            currentSize = 0
        } 
        fileCount += 1 
        currentSize = currentSize + size
        m3uSplitFileContents[m3uSplitFileContentsIndex] += v
        m3uSplitFileContents[m3uSplitFileContentsIndex] += "\n"

    }

    //Write the files in same directory as origin m3u.
    for i := 0; i < len(m3uSplitFileContents); i++ {
        newM3uBasename := filepath.Base(m3uFilename) + "-" + strconv.Itoa(i) + ".m3u"
        newM3uFilename := filepath.Join(m3uFilenameDir,newM3uBasename)
        dj80hdutil.StringToFile(newM3uFilename,m3uSplitFileContents[i])
        fmt.Printf("CREATED FILE %q\n",newM3uFilename)
    }

}

