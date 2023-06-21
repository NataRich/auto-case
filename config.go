package main

import (
  "os"
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
)

// 申请人或被申请人设置
type PersonConfig struct {
  // 当事人类型
  Type                string          `json:"type"`

  // 姓名
  Name                string          `json:"name"`

  // 手机号码
  Tel                 string          `json:"tel"`

  // 证件类型
  CredentialsType     string          `json:"credentialsType"`

  // 证件号码
  IDCardNo            string          `json:"idCardNo"`

  // 性别
  Sex                 string          `json:"sex"`

  // 出生日期
  Birthday            string          `json:"birthday"`

  // 民族
  Nation              string          `json:"nation"`

  // 居住地代码
  AreaCode            string          `json:"areaCode"`

  // 居住地址
  Address             string          `json:"address"`
}

// 案件设置
type CaseConfig struct {
  // 调解类型
  Type                string          `json:"type"`

  // 案件年份
  Year                string          `json:"year"`

  // 案件类型
  CaseCatalog         string          `json:"caseCatalog"`

  // 纠纷类型
  DisputeType         string          `json:"disputeType"`

  // 案由
  CauseCode           string          `json:"causeCode"`

  // 案件状态
  State               string          `json:"state"`

  // 成功状态
  SuccessState        string          `json:"successState"`

  // 调解开始日期
  StartTime           string          `json:"startTime"`

  // 调解结束日期
  EndTime             string          `json:"endTime"`

  // 纠纷概况
  Dispute             string          `json:"dispute"`

  // 调解方案
  Agreement           string          `json:"agreement"`

  // 自动生成调解协议
  AutoCreate          string          `json:"autoCreate"`

  // 默认申请人信息
  DefaultApplicant    *PersonConfig   `json:"defaultApplicant"`

  // 默认被申请人信息
  DefaultRespondent   *PersonConfig   `json:"defaultRespondent"`

  // 默认调解员ID
  DefaultMediatorId   string          `json:"defaultMediatorId"`
}

// 数据源配置
type DataConfig struct {
  // 数据Excel表
  Path          string                `json:"path"`

  // 工作表名
  Sheet         string                `json:"sheet"`

  // 跳过列名行
  SkipHeader    bool                  `json:"skipHeader"`

  // 跳过行数
  SkipLines     int                   `json:"skipLines"`

  // 截止行数
  ExecCount     int                   `json:"execCount"`

  // 申请人列号
  ApplicantCol  string                `json:"applicantCol"`

  // 被申请人列号
  RespondentCol string                `json:"respondentCol"`

  // 自定义配置
  Mapper        map[string]string     `json:"mapper"`
}

// 请求配置
type RequestConfig struct {
  // 单次请求延迟（秒）
  Delay       int             `json:"delay"`

  // 单次请求重试次数
  Retry       int             `json:"retry"`

  // 单次请求超时时长（秒）
  Timeout     int             `json:"timeout"`

  // 浏览器中的完整cookie字符串
  Cookie      string          `json:"cookie"`
}

// 调试配置
type DebugConfig struct {
  // 调试模式
  Verbose     bool            `json:"verbose"`

  // 伪请求模式
  Fake        bool            `json:"fake"`

  // 错误日志文件
  LogPath     string          `json:"logPath"`
}

// 全局配置
type GlobalConfig struct {
  Case        *CaseConfig     `json:"case"`
  Data        *DataConfig     `json:"data"`
  Request     *RequestConfig  `json:"request"`
  Debug       *DebugConfig    `json:"debug"`
}

var Conf *GlobalConfig

// 默认配置
func InitConf() (*GlobalConfig) {
  return &GlobalConfig{
    Case:     &CaseConfig{
      Type:               "0",
      Year:               NowYearStr(),
      CaseCatalog:        "",
      DisputeType:        "",
      CauseCode:          "",
      State:              "",
      SuccessState:       "",
      StartTime:          "",
      EndTime:            "",
      Dispute:            "",
      Agreement:          "",
      AutoCreate:         "1",
      DefaultMediatorId:  "",
      DefaultApplicant:   &PersonConfig{
        Type:             "",
        Name:             "",
        Tel:              "",
        CredentialsType:  "",
        IDCardNo:         "",
        Sex:              "",
        Birthday:         "",
        Nation:           "",
        AreaCode:         "",
        Address:          "",
      },
      DefaultRespondent:  &PersonConfig{
        Type:             "",
        Name:             "",
        Tel:              "",
        CredentialsType:  "",
        IDCardNo:         "",
        Sex:              "",
        Birthday:         "",
        Nation:           "",
        AreaCode:         "",
        Address:          "",
      },
    },

    Data:     &DataConfig{
      Path:               "",
      Sheet:              "",
      SkipHeader:         false,
      SkipLines:          0,
      ExecCount:          1,
      ApplicantCol:       "",
      RespondentCol:      "",
      Mapper:             map[string]string{},
    },

    Request:  &RequestConfig{
      Delay:              2,
      Retry:              3,
      Timeout:            10,
      Cookie:             "",
    },

    Debug:    &DebugConfig{
      Verbose:            true,
      Fake:               false,
      LogPath:            "error.log",
    },
  }
}

// 从文件加载配置
func LoadConf(path string) (*GlobalConfig, error) {
  if _, err := os.Stat(path); err != nil {
    return nil, err
  }

  bytes, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  if err := json.Unmarshal(bytes, &Conf); err != nil {
    return nil, err
  }

  return Conf, nil
}

// 更新当事人姓名
func UpdateNames(appName string, resName string) error {
  if Conf == nil {
    return fmt.Errorf("全局配置为空\n")
  }

  if Conf.Case == nil {
    return fmt.Errorf("案件配置为空\n")
  }

  if Conf.Case.DefaultApplicant == nil {
    return fmt.Errorf("默认申请人为空\n")
  } else {
    Conf.Case.DefaultApplicant.Name = appName
  }

  if Conf.Case.DefaultRespondent == nil {
    return fmt.Errorf("默认被申请人为空\n")
  } else {
    Conf.Case.DefaultRespondent.Name = resName
  }

  return nil
}

// 保存配置到文件
func SaveConf(path string) error {
  if Conf == nil {
    return fmt.Errorf("没有加载配置，无法保存\n")
  }

  data, err := json.MarshalIndent(Conf, "", "  ")
  if err != nil {
    return err
  }

  if err := ioutil.WriteFile(path, data, 0644); err != nil {
    return err
  }

  return nil
}

// 预先配置检查
func PreCheck(conf *GlobalConfig) bool {
  if conf == nil {
    log.Println("全区配置为空")
    return false
  }

  // 案件配置检查
  ca := conf.Case
  if ca == nil {
    log.Println("案件配置不得为空")
    return false
  }

  if ca.Type == "" {
    log.Println("调解类型不得为空")
    return false
  }

  if ca.Year == "" {
    log.Println("案件年份不得为空")
    return false
  }

  if ca.CaseCatalog == "" {
    log.Println("案件类型不得为空")
    return false
  }

  if ca.DisputeType == "" {
    log.Println("纠纷类型不得为空")
    return false
  }

  if ca.CauseCode == "" {
    log.Println("案由不得为空")
    return false
  }

  if ca.State == "" {
    log.Println("案件状态不得为空")
    return false
  }

  if ca.SuccessState == "" {
    log.Println("成功状态为空（允许，但请注意是否符合表单要求）")
    return false
  }

  if ca.StartTime == "" {
    log.Println("调解开始日期为空（允许，将被随机日期覆写）")
  } else {
    log.Println("调解开始日期不为空（允许，将被随机日期覆写）")
  }

  if ca.EndTime == "" {
    log.Println("调解结束日期为空（允许，将被随机日期覆写）")
  } else {
    log.Println("调解结束日期不为空（允许，将被随机日期覆写）")
  }

  if ca.Dispute == "" {
    log.Println("纠纷概况不得为空")
    return false
  }

  if ca.Agreement == "" {
    log.Println("调解方案不得为空")
    return false
  }

  if ca.AutoCreate == "" {
    log.Println("是否自动生成调解协议不得为空（是：1；否：0）")
    return false
  }

  if ca.DefaultMediatorId == "" {
    log.Println("默认调解员ID不得为空")
    return false
  }

  if ca.DefaultApplicant == nil {
    log.Println("默认申请人信息为空")
    return false
  }

  if ca.DefaultRespondent == nil {
    log.Println("默认被申请人信息为空")
    return false
  }


  // 数据源配置检查
  data := conf.Data
  if data == nil {
    log.Println("数据源配置不得为空")
    return false
  }

  if data.Path == "" {
    log.Println("excel数据表路径不得为空")
    return false
  }

  if data.Sheet == "" {
    log.Println("excel工作表名不得为空")
    return false
  }

  if data.SkipLines < 0 {
    log.Println("跳过行数不得为负数")
    return false
  }

  if data.ExecCount <= 0 {
    log.Println("执行行数不得为0或负数")
    return false
  }

  if data.ApplicantCol == "" {
    log.Println("申请人列号不得为空")
    return false
  }

  if data.RespondentCol == "" {
    log.Println("被申请人列号不得为空")
    return false
  }

  // 请求配置检查
  req := conf.Request
  if req == nil {
    log.Println("请求配置不得为空")
    return false
  }

  if req.Delay < 0 {
    log.Println("单次请求延迟不得为负数")
    return false
  }

  if req.Retry < 0 {
    log.Println("单次请求重试次数不得为负数")
    return false
  }

  if req.Timeout < 0 {
    log.Println("单次请求超时时长不得为负数")
    return false
  }

  if req.Cookie == "" {
    log.Println("请求Cookie配置不得为空")
    return false
  }


  // 调试配置检查
  debug := conf.Debug
  if debug == nil {
    log.Println("调试配置为空")
    return false
  }

  if debug.LogPath == "" {
    log.Println("错误日志路径为空（不必要，但强烈建议配置！）")
  }
  
  return true
}

// 当事人配置检查
func PersonCheck(per *PersonConfig) bool {
  if per == nil {
    log.Println("申请人/被申请人信息为空")
    return false
  }

  if per.Type == "" {
    log.Println("当事人类型为空")
    return false
  }

  if per.Name == "" {
    log.Println("当事人姓名为空")
    return false
  }

  if per.Tel == "" {
    log.Println("当事人手机号为空")
    return false
  }

  if per.CredentialsType == "" {
    log.Println("当事人证件类型为空")
    return false
  }

  if per.IDCardNo == "" {
    log.Println("当事人身份证号为空（允许，但请检查是否符合需求）")
  }

  if per.Sex == "" {
    log.Println("当事人性别为空")
    return false
  }

  if per.Birthday == "" {
    log.Println("当事人生日为空")
    return false
  }

  if per.Nation == "" {
    log.Println("当事人民族为空")
    return false
  }

  if per.AreaCode == "" {
    log.Println("当事人地区代号为空")
    return false
  }

  if per.Address == "" {
    log.Println("当事人地址为空")
    return false
  }

  return true
}

// debug 打印
func DebugPrint(msg string) {
  if Conf != nil && Conf.Debug.Verbose {
    fmt.Println(msg)
  }
}
