package logging

import "log"

func CmpInfo(disCode string, v ...interface{}) {
	var err error
	filePath := getLogFilePath()
	fileName := disCode + getLogFileName()

	Fcmp, err := openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}
	loggercpm := log.New(Fcmp, DefaultPrefix, log.LstdFlags)
	setPrefix(INFO)
	loggercpm.Println(v)

}
