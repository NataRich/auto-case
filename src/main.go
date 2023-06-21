package main

import (
  "os"
  "fmt"
  "log"

  "github.com/urfave/cli/v2"
  "github.com/xuri/excelize/v2"
)

const (
  CONFIG_FILE = "config.json"
)

// 初始化默认配置
func initCase(ctx *cli.Context) error {
  Conf = InitConf()
  return SaveConf(CONFIG_FILE)
}

// 调用接口新建案例
func newCase(ctx *cli.Context) error {
  if _, err := LoadConf(CONFIG_FILE); err != nil {
    return err
  }

  if !PreCheck(Conf) {
    return fmt.Errorf("预先检查失败，请检查配置文件\n")
  }

  f, err := excelize.OpenFile(Conf.Data.Path)
  if err != nil {
    return err
  }
  
  defer func() {
    if err := f.Close(); err != nil {
      log.Println(err)
    }
  }()

  rows, err := f.Rows(Conf.Data.Sheet)
  if err != nil {
    return err
  }

  baseLine := 0
  if Conf.Data.SkipHeader {
    rows.Next()
    baseLine += 1
  }

  for i := 0; i < Conf.Data.SkipLines; i++ {
    rows.Next();
  }

  baseLine += Conf.Data.SkipLines
  DebugPrint(fmt.Sprintf("跳过excel表的前%d行（%s）", Conf.Data.SkipLines, Conf.Data.Path))

  count := Conf.Data.ExecCount
  for rows.Next() {
    if count <= 0 {
      break
    }

    count--

    InsertRandomDates(Conf.Case)

    DebugPrint(fmt.Sprintf("正在抓取excel表的第%d行数据（%s）", 
                            baseLine + Conf.Data.ExecCount - count,
                            Conf.Data.Path))

    row, err := rows.Columns()
    if err != nil {
      DebugPrint("无法获取该行内容，将跳过该行")
      LogError(err, Conf.Debug.LogPath)
      continue
    }

    appCol, err := excelize.ColumnNameToNumber(Conf.Data.ApplicantCol)
    if err != nil {
      DebugPrint(fmt.Sprintf("无法获取该行%s列名对应的索引，将跳过该行", 
                              Conf.Data.ApplicantCol))
      LogError(err, Conf.Debug.LogPath)
      continue
    }

    resCol, err := excelize.ColumnNameToNumber(Conf.Data.RespondentCol)
    if err != nil {
      DebugPrint(fmt.Sprintf("无法获取该行%s列名对应的索引，将跳过该行", 
                              Conf.Data.RespondentCol))
      LogError(err, Conf.Debug.LogPath)
      continue
    }
    
    appName := row[appCol - 1]
    resName := row[resCol - 1]

    DebugPrint(fmt.Sprintf("申请人：%s", appName))
    DebugPrint(fmt.Sprintf("被申请人：%s", resName))

    if err := UpdateNames(appName, resName); err != nil {
      DebugPrint("更新申请人姓名失败，请检查配置文件/数据源后重试")
      DebugPrint(fmt.Sprintf("%v", err.Error()))
      LogError(err, Conf.Debug.LogPath)
      return err
    }

    if !PersonCheck(Conf.Case.DefaultApplicant) {
      DebugPrint("申请人信息检查失败，将跳过该行")
      continue
    }

    if !PersonCheck(Conf.Case.DefaultRespondent) {
      DebugPrint("申请人信息检查失败，将跳过该行")
      continue
    }

    if err := MakeRequestWithRetry(Conf.Case, Conf.Request, Conf.Debug.Fake); err != nil {
      DebugPrint(fmt.Sprintf("重试了%d次，新建请求仍旧失败，将跳过该行",
                              Conf.Request.Retry))
      LogError(err, Conf.Debug.LogPath)
    }
  }

  return nil
}

func main() {
  cli.CommandHelpTemplate = `程序:
   {{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}
使用:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
命令:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
选项:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

  app := cli.NewApp()
	app.Usage = "人民法院调解新建案例接口"
	app.UsageText = "case COMMANDS [ARG...]"
	app.ArgsUsage = "ArgsUsage"
	app.EnableBashCompletion = true
	app.HideVersion = true
  app.Commands = []*cli.Command{
    &cli.Command{
      Name: "init",
      Usage: "Initializes a configuration file",
      UsageText: "case init",
      Action: initCase,
    },
    &cli.Command{
      Name: "new",
      Usage: "Creates a new case record",
      UsageText: "case new [OPTIONS...]",
      Flags: []cli.Flag{
        &cli.StringFlag{
          Name:     "output",
          Aliases:  []string{"o"},
          Usage:    "Writes response to the given file",
        },
      },
      Action: newCase,
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
