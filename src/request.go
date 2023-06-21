package main

import (
  "io"
  "fmt"
  "log"
  "time"
  "strings"
  "net/url"
  "net/http"
  "encoding/json"
)


type ApplicantBody struct {
  Type string             `json:"applicantType"`
  Name string             `json:"applicantName"`
  Tel string              `json:"applicantTel"`
  CredentialsType string  `json:"credentialsType"`
  CredentialsName string  `json:"credentialsName"`
  IDCardNo string         `json:"applicantIDCardNo"`
  Sex string              `json:"applicantSex"`
  Birthday string         `json:"applicantBirthday"`
  Nation string           `json:"applicantNation"`
  AreaCode string         `json:"areaCode"`
  Address string          `json:"applicantAddress"`
  Email string            `json:"email"`
  AgentList []string      `json:"agentList"`
  FileList []string       `json:"fileList"`
}

type RespondentBody struct {
  Type string             `json:"respondentType"`
  Name string             `json:"respondentName"`
  Tel string              `json:"respondentTel"`
  StaticPhone string      `json:"respondentStaticPhone"`
  CredentialsType string  `json:"credentialsType"`
  CredentialsName string  `json:"credentialsName"`
  IDCardNo string         `json:"respondentIDCardNo"`
  Sex string              `json:"respondentSex"`
  Birthday string         `json:"respondentBirthday"`
  Nation string           `json:"respondentNation"`
  AreaCode string         `json:"areaCode"`
  Address string          `json:"respondentAddress"`
  Email string            `json:"email"`
  AgentList []string      `json:"agentList"`
  FileList []string       `json:"fileList"`
}


type CaseBody struct {
  Type string             `json:"type"`
  Year string             `json:"year"`
  DraftFlag string        `json:"draftFlag"`
  CaseCatalog string      `json:"caseCatalog"`
  DisputeType string      `json:"disputeType"`
  CauseCode string        `json:"causeCode"`
  MediationCaseNo string  `json:"mediationCaseNo"`
  Money string            `json:"money"`
  ClaimMoney string       `json:"claimMoney"`
  State string            `json:"state"`
  SuccessState string     `json:"successState"`
  Remark string           `json:"remark"`
  StartTime string        `json:"startTimeStr"`
  EndTime string          `json:"endTimeStr"`
  Dispute string          `json:"dispute"`
  Agreement string        `json:"agreement"`
  MediatorId string       `json:"mediatorId"`
  AutoCreate string       `json:"autoCreate"`
  DocList []string        `json:"docList"`
  NoteList []string       `json:"noteList"`

  ApplicantList  []*ApplicantBody  `json:"applicantPartyList"`
  RespondentList []*RespondentBody `json:"respondentPartyList"`

  Evidences []string      `json:"evidences"`
}

var appBody = &ApplicantBody{
  Type:             "",
  Name:             "",
  Tel:              "",
  CredentialsType:  "",
  CredentialsName:  "",
  IDCardNo:         "",
  Sex:              "",
  Birthday:         "",
  Nation:           "",
  AreaCode:         "",
  Address:          "",
  Email:            "",
  AgentList:        []string{},
  FileList:         []string{},
}

var resBody = &RespondentBody{
  Type:             "",
  Name:             "",
  Tel:              "",
  StaticPhone:      "",
  CredentialsType:  "",
  CredentialsName:  "",
  IDCardNo:         "",
  Sex:              "",
  Birthday:         "",
  Nation:           "",
  AreaCode:         "",
  Address:          "",
  Email:            "",
  AgentList:        []string{},
  FileList:         []string{},
} 


var caseBody = &CaseBody{
  Type:             "",
  Year:             "",
  DraftFlag:        "1",
  CaseCatalog:      "",
  DisputeType:      "",
  CauseCode:        "",
  MediationCaseNo:  "",
  Money:            "",
  ClaimMoney:       "",
  State:            "",
  SuccessState:     "",
  Remark:           "",
  StartTime:        "",
  EndTime:          "",
  Dispute:          "",
  Agreement:        "",
  MediatorId:       "",
  AutoCreate:       "1",
  DocList:          []string{},
  NoteList:         []string{},

  ApplicantList:    []*ApplicantBody{},
  RespondentList:   []*RespondentBody{},

  Evidences:        []string{},
}

const ENDPOINT string = "http://tiaojie.court.gov.cn/fayuan/a/offline/addOffline"

func setAppBody(conf *PersonConfig) {
  appBody.Type            = conf.Type
  appBody.Name            = conf.Name
  appBody.Tel             = conf.Tel
  appBody.CredentialsType = conf.CredentialsType
  appBody.IDCardNo        = conf.IDCardNo
  appBody.Sex             = conf.Sex
  appBody.Birthday        = conf.Birthday
  appBody.AreaCode        = conf.AreaCode
  appBody.Address         = conf.Address
}

func setResBody(conf *PersonConfig) {
  resBody.Type            = conf.Type
  resBody.Name            = conf.Name
  resBody.Tel             = conf.Tel
  resBody.CredentialsType = conf.CredentialsType
  resBody.IDCardNo        = conf.IDCardNo
  resBody.Sex             = conf.Sex
  resBody.Birthday        = conf.Birthday
  resBody.Nation          = conf.Nation
  resBody.AreaCode        = conf.AreaCode
  resBody.Address         = conf.Address
}

func setBody(ca *CaseConfig) {
  caseBody.Type           = ca.Type
  caseBody.Year           = ca.Year
  caseBody.CaseCatalog    = ca.CaseCatalog
  caseBody.DisputeType    = ca.DisputeType
  caseBody.CauseCode      = ca.CauseCode
  caseBody.State          = ca.State
  caseBody.SuccessState   = ca.SuccessState
  caseBody.StartTime      = ca.StartTime
  caseBody.EndTime        = ca.EndTime
  caseBody.Dispute        = ca.Dispute
  caseBody.Agreement      = ca.Agreement
  caseBody.AutoCreate     = ca.AutoCreate
  caseBody.MediatorId     = ca.DefaultMediatorId

  setAppBody(ca.DefaultApplicant)
  setResBody(ca.DefaultRespondent)

  caseBody.ApplicantList  = []*ApplicantBody{ appBody }
  caseBody.RespondentList = []*RespondentBody{ resBody }

  DebugPrint(fmt.Sprintf("已成功载入案件配置"))
}


func MakeRequest(body *CaseBody, cookie string, timeout int, fake bool) error {
  if body == nil || cookie == "" || timeout < 0 {
    return fmt.Errorf("MakeRequest()参数错误")
  }

  // body序列化
  s, err := json.Marshal(body)
  if err != nil {
    DebugPrint("请求时，序列化错误")
    return err
  }

  // 设置表单数据
  form := url.Values{
    "mediationFormStr": { string(s) },
  }

  request, err := http.NewRequest("POST", ENDPOINT, strings.NewReader(form.Encode()))
  if err != nil {
    DebugPrint("请求时，请求创建错误")
    return err
  }

  // 设置请求头
  request.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
  request.Header.Add("Accept-Encoding", "gzip, deflate, br")
  request.Header.Add("Connection", "keep-alive")
  request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  request.Header.Add("Cookie", cookie)
  request.Header.Add("Host", "tiaojie.court.gov.cn")
  request.Header.Add("Origin", "http://tiaojie.court.gov.cn")
  request.Header.Add("Referer", "http://tiaojie.court.gov.cn/fayuan/offline/toAddOffline")
  request.Header.Add("X-Requested-With", "XMLHttpRequest")
  request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/114.0")

  // 若为伪请求模式，打印所有相关信息
  if fake {
    DebugPrint(fmt.Sprintf("接口 （POST）: %s", ENDPOINT))

    DebugPrint("请求头：")
    for name, values := range request.Header {
      for _, value := range values {
        DebugPrint(fmt.Sprintf("%s: %s", name, value))
      }
    }

    DebugPrint("请求体：")
    xx, err := json.MarshalIndent(body, "", "  ")
    if err != nil {
      DebugPrint("无法序列化")
      return err
    }
    DebugPrint(fmt.Sprintf("%s", string(xx)))

    DebugPrint("请求体（表单格式）：")
    DebugPrint(fmt.Sprintf("%s", string(s)))
    return nil
  }

  // 否则，发送真实请求
  client := http.Client{
    Timeout: time.Duration(timeout) * time.Second,
  }

  response, err := client.Do(request)
  if err != nil {
    DebugPrint(fmt.Sprintf("无法发送请求"))
    return err
  }

  // 判断返回体
  if response.StatusCode != 200 {
    return fmt.Errorf("新建失败，返回值：%d", response.StatusCode)
  }

  defer response.Body.Close()
  bytes, err := io.ReadAll(response.Body)
  if err != nil {
    log.Fatal("新建结果未知，无法解析响应请求体，请手动确认\n")
  }

  rbody := string(bytes)
  if strings.Contains(rbody, "html") {
    log.Fatal("新建失败，请立即更新Cookie信息\n")
  } 

  if strings.Contains(rbody, "-1") {
    return fmt.Errorf("新建失败，返回码为-1：%s", rbody)
  }

  fmt.Println("新建成功！")
  fmt.Println(rbody)
  return nil
}


func MakeRequestWithRetry(caseConf *CaseConfig, reqConf *RequestConfig, fake bool) error {
  if reqConf == nil {
    return fmt.Errorf("请求配置不能为空")
  }

  if reqConf.Cookie == "" {
    return fmt.Errorf("Cookie中的acw_tc参数不能为空（请自行到浏览器登录后复制）")
  }

  if caseConf == nil {
    return fmt.Errorf("案件配置不能为空")
  }

  // 加载默认参数
  setBody(caseConf)
  
  // 发送请求
  if err := MakeRequest(caseBody, reqConf.Cookie, reqConf.Timeout, fake); err != nil {
    DebugPrint(fmt.Sprintf("首次请求失败，即将重试（预计重试%d次）", 
                            reqConf.Retry))

    var err error
    for i := 0; i < reqConf.Retry; i++ {
      time.Sleep(time.Duration(reqConf.Delay) * time.Second)
      DebugPrint(fmt.Sprintf("已等待%d秒，再次尝试", reqConf.Delay))

      err = MakeRequest(caseBody, reqConf.Cookie, reqConf.Timeout, fake)
      if err != nil {
        DebugPrint(fmt.Sprintf("第%d次重试失败", i + 1))
      }
    }

    if err != nil {
      return err
    }
  } else {
    time.Sleep(time.Duration(reqConf.Delay) * time.Second)
    DebugPrint("休息一下...\n")
  }

  return nil
}
