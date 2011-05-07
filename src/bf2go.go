package bf2go

import(
	"os"
	"log"
	"strconv"
)

type brainfuck struct {
	SourceFile *os.File
	DestFile *os.File
	Token []byte
	Line int
	Char int
	Debug bool
}

func Translate(source string, dest string, debug bool) {

	in,_ := os.Open(source)
	out,_ := os.OpenFile(dest,os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)

	bf := &brainfuck{in,out,make([]byte,1),1,1,debug}

	bf.SourceFile.Seek(0,0)
	bf.DestFile.Seek(0,0)

	bf.initalizeOutputFile()

	bf.ParseToken()
}

func (bf *brainfuck) initalizeOutputFile() {
	bf.DestFile.WriteString("package main\n")
	bf.DestFile.WriteString("import (\n")
	bf.DestFile.WriteString("\t\"os\"\n")
	bf.DestFile.WriteString("\t\"bufio\"\n")
	if bf.Debug {	
		bf.DestFile.WriteString("\t\"log\"\n")
	}
	bf.DestFile.WriteString(")\n")
	bf.DestFile.WriteString("\n\n");

	bf.DestFile.WriteString("func coredump(stack []byte){\n")
	bf.DestFile.WriteString("\tfile,_ := os.OpenFile(\"bfcoredump\",os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)\n")
	bf.DestFile.WriteString("\tfile.Write(stack)\n")
	bf.DestFile.WriteString("\tfile.Close()\n")	
	bf.DestFile.WriteString("}\n")

	bf.DestFile.WriteString("\n\n");
	bf.DestFile.WriteString("func main(){\n")
	bf.DestFile.WriteString("\tstack := make([]byte,30000)\n")
	bf.DestFile.WriteString("\tdefer func() {\n\tcoredump(stack)\n\t}()\n")
	bf.DestFile.WriteString("\tfor i:=0;i<len(stack);i++ {\nstack[i]=0\n}\n")
	bf.DestFile.WriteString("\tstackPosition := 0\n")
	bf.DestFile.WriteString("\twriter := bufio.NewWriter(os.Stdout)\n")
	bf.DestFile.WriteString("\treader := bufio.NewReader(os.Stdin)\n")
	bf.DestFile.WriteString("\tif writer == nil{\nwriter=nil}\n")
	bf.DestFile.WriteString("\tif reader == nil{\nreader=nil}\n")

	bf.DestFile.WriteString("\tbuf := make([]byte,1)\n")
	bf.DestFile.WriteString("\tif buf == nil{\nbuf=nil}\n")

}

func (bf *brainfuck) finnishOutputFile() {
	bf.DestFile.WriteString("\twriter.Flush()\n")
	bf.DestFile.WriteString("}\n")
}

func (bf *brainfuck) ParseToken() {
	if _,err := bf.SourceFile.Read(bf.Token); err != nil {
		bf.finnishOutputFile();
		log.Printf("%s",err)
		return
	}

	token := bf.Token[0];

	if token == '\n' {
		bf.Line++
		bf.Char=1
	} else if token != '\r' {
		bf.Char++
	}

	switch(token) {
		case '<':
			bf.DestFile.WriteString("\n\tstackPosition--\n\tif stackPosition < 0 {\n\tstackPosition=0\n\t}\n")
		break;
		case '>':
			bf.DestFile.WriteString("\n\tstackPosition++\n")
			bf.DestFile.WriteString("\tif stackPosition >= len(stack) {\n")
			if bf.Debug {			
				bf.DestFile.WriteString("\tlog.Printf(\"increasing stack length to %d\",len(stack)+2)\n")
			}
			bf.DestFile.WriteString("\ttmp := make([]byte,len(stack)+2)\n")
			bf.DestFile.WriteString("\tfor i:= 0;i<len(stack);i++ {\n\ttmp[i]=stack[i]\n\t}\n")
			bf.DestFile.WriteString("\tstack=tmp\n")
			bf.DestFile.WriteString("\t}\n")
		break;
		case '+':
			bf.DestFile.WriteString("\tstack[stackPosition]++\n")
		break;
		case '-':
			bf.DestFile.WriteString("\tstack[stackPosition]--\n")
		break;
		case '.':
			if bf.Debug {
				bf.DestFile.WriteString("\tlog.Printf(\"print@%d\",stackPosition)\n")
			}
			bf.DestFile.WriteString("\twriter.WriteByte(stack[stackPosition])\n")
		break;
		case ',':
			
			//bf.DestFile.WriteString("if num,_ := os.Stdin.Read(buf); num > 0 {\nstack[stackPosition]=buf[0]\n}\n")
			bf.DestFile.WriteString("\tstack[stackPosition],_=reader.ReadByte()\n")
		break;
		case '|':
			bf.DestFile.WriteString("\tswitch(stack[stackPosition]) {\n")
			bf.DestFile.WriteString("\tcase 1:\n")
			bf.DestFile.WriteString("\tend := 0\n")
			bf.DestFile.WriteString("\tstart := stackPosition+1\n")
			bf.DestFile.WriteString("\tfor i:=start;i<len(stack);i++ {\n")
			bf.DestFile.WriteString("\tif stack[i] != 0 {\n\tend++\n\t} else {\n\tbreak;\n\t}\n")
			bf.DestFile.WriteString("\t}\n")
			bf.DestFile.WriteString("\tfile,error := os.OpenFile((string)(stack[start:end]),os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)\n")
			bf.DestFile.WriteString("\tif error == nil { writer = bufio.NewWriter(file) } else { writer.Write(([]byte)(error.String())) }\n")
			bf.DestFile.WriteString("\tbreak;\n")
			bf.DestFile.WriteString("\tdefault:\n")
			bf.DestFile.WriteString("\twriter = bufio.NewWriter(os.Stdout)\n")
			bf.DestFile.WriteString("\tbreak;\n")
			bf.DestFile.WriteString("\t}\n")

		break;
		case '[':
			if bf.Debug {
				bf.DestFile.WriteString("\tlog.Printf(\"loop start@"+strconv.Itoa(bf.Line)+":"+strconv.Itoa(bf.Char)+"\")\n")
			}
			bf.DestFile.WriteString("\tfor{\n")
		break;
		case ']':
			if bf.Debug {
				bf.DestFile.WriteString("\tlog.Printf(\"loop end@"+strconv.Itoa(bf.Line)+":"+strconv.Itoa(bf.Char)+", stack=%d,stackPosition=%d\",stack[stackPosition],stackPosition)\n")
			}
			bf.DestFile.WriteString("\tif stack[stackPosition] == 0 {\n\t\tbreak;\n\t}\n")
			bf.DestFile.WriteString("\t}\n")
		break;
		case '}':
			bf.DestFile.WriteString("\tgo func(stackPosition int){\n")
		break;
		case '{':
			bf.DestFile.WriteString("\t}(stackPosition+0)\n")
		break;
		default:
			log.Printf("ignore '%c'",token)
		break;
	}

	bf.ParseToken()
}
