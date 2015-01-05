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

//HELP Text displayed when no arguments.
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
    //FIXME - Check valid m3u file and valid split length

    //First argument is the name of the m3u to be split (e.g. foo/bar/baz.m3u)
    m3uFilename := args[0]

    //Split Lenth is parsed into an int and assumed to be in Megabytes
    m3uSplitLength,_ := strconv.ParseInt(args[1],0,64)
    m3uSplitLength = m3uSplitLength * 1024 * 1024 //700 becomes 700M

    //Data Structure to hold m3u files that are split
    //FIXME this could be done so much better.
    var m3uSplitFileContents []string
    m3uSplitFileContents = append(m3uSplitFileContents,"")
    m3uSplitFileContentsIndex := 0

    //Total Count of files referenced in the original m3u file
    fileCount := 0

    //We must know the directory containing the m3u file because its 
    //contents likely have filenames relative to that directory.
    m3uFilenameDir := filepath.Dir(m3uFilename)

    //Get the entire contents of this m3u file as an array of strings
    mp3Files := dj80hdutil.FileToArray(m3uFilename)

    //currentSize accumualtes the size as we move through each file
    //so when we know when to split.
    currentSize := int64(0)        

    //Loop for every path in the m3u file
    for i := 0; i < len(mp3Files); i++ {
        //Get the relative path in the m3u file
        v := mp3Files[i]

        //ALSO resolve the name to the original m3u file
        //because the directory where this is running is likely NOT
        //the same directory as the m3u file.
        resolvedFilename := filepath.Join(m3uFilenameDir,v) 

        //Get the size of the file mentioned in the m3u file
        fi,err := os.Stat(resolvedFilename)
        if err != nil {
            fmt.Printf("WARNING: Could not stat file %q\n",v)
            continue
        }
        size := fi.Size()

        //If we are at a point where we need to split, take care of that.
        if (currentSize + size) > m3uSplitLength {
            //Create a new String to accumulate file contents
            m3uSplitFileContents = append(m3uSplitFileContents,"")

            //Update the index so we can append to the correct string
            m3uSplitFileContentsIndex = m3uSplitFileContentsIndex + 1

            //Reset the current size.
            currentSize = 0
        } 

        //Total files processed always is incremented.
        fileCount += 1 

        //Increment our currentSize by the size of the file referenced in m3u
        currentSize = currentSize + size

        //Add the UNRESOLVED path to the file contents.
        //(UNRESOLVED because new m3u files are placed alongside the original)
        m3uSplitFileContents[m3uSplitFileContentsIndex] += v
        m3uSplitFileContents[m3uSplitFileContentsIndex] += "\n"

    }

    //Convert our data structure to actual new m3u files. 
    for i := 0; i < len(m3uSplitFileContents); i++ {
        newM3uBasename := filepath.Base(m3uFilename) + "-" + strconv.Itoa(i) + ".m3u"
        newM3uFilename := filepath.Join(m3uFilenameDir,newM3uBasename)
        dj80hdutil.StringToFile(newM3uFilename,m3uSplitFileContents[i])
        fmt.Printf("CREATED FILE %q\n",newM3uFilename)
    }

}

