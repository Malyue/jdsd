package jdsd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// request url
var request_url = "https://jdsd.gzhu.edu.cn/coctl_gzhu/index_wx.php"

// request header
var userAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF XWEB/6398"
var Host = "jdsd.gzhu.edu.cn"
var Accept = "*/*"
var AcceptEncoding = "gzip,deflate,br"
var AcceptLanguage = "zh-CN,zh"

// addHeader add the necessary header for the request
func addHeader(request *http.Request) *http.Request {
	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Host", Host)
	request.Header.Add("Accept", Accept)
	request.Header.Add("Accept-Encoding", AcceptEncoding)
	request.Header.Add("Accept-Language", AcceptLanguage)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

// GetUserInfo Get the user info,includes credits,per_day_credits,user_name,xh
func GetUserInfo(key string) (error, map[string]interface{}) {
	client := &http.Client{}
	//data := strings.NewReader("route=user_info&key=" + key)
	// construct the request body
	data := url.Values{"route": {"user_info"}, "key": {key}}
	body := strings.NewReader(data.Encode())
	request, err := http.NewRequest("POST", request_url, body)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	request = addHeader(request)
	resp, err := client.Do(request)
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	info := make(map[string]interface{})
	err = json.Unmarshal(resp_body, &info)
	if err != nil {
		return err, nil
	}
	return nil, info
}

// Per_day_question Finish the question task,have three times in one day
func Per_day_question(key string) error {
	client := &http.Client{}
	// get the question list
	data := url.Values{"route": {"train_list_get"}, "diff": {"0"}, "key": {key}}
	body := strings.NewReader(data.Encode())
	request_question, err := http.NewRequest("POST", request_url, body)

	if err != nil {
		fmt.Println(err)
		return err
	}
	request_question = addHeader(request_question)
	resp_question_list, err := client.Do(request_question)
	defer resp_question_list.Body.Close()

	resp_question_list_body, err := io.ReadAll(resp_question_list.Body)
	question_list := make(map[string]interface{})
	err = json.Unmarshal(resp_question_list_body, &question_list)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// the body we need is question_list["re"]["question_bag"][index]["num"]
	// so we should parse the struct
	re := question_list["re"].(map[string]interface{})
	question_bag := re["question_bag"].([]interface{})
	train_result := make([][2]string, 0)
	for _, question := range question_bag {
		item := [2]string{question.(map[string]interface{})["num"].(string), "1"}
		train_result = append(train_result, item)
	}
	train_id := re["train_id"].(string)
	train_result_str, err := json.Marshal(train_result)

	// commit the request of complete
	data = url.Values{"route": {"train_finish"}, "train_id": {train_id}, "train_result": {string(train_result_str)}, "key": {key}}
	body = strings.NewReader(data.Encode())
	request_finish, err := http.NewRequest("POST", request_url, body)
	if err != nil {
		fmt.Println()
		return err
	}
	request_finish = addHeader(request_finish)
	resp, err := client.Do(request_finish)
	defer resp.Body.Close()
	return err
}

func Read(key string) error {
	client := &http.Client{}
	for read_type := 1; read_type < 6; read_type++ {
		// commit begin request
		data := url.Values{"route": {"classic_time"}, "addtime": {"0"}, "type": {strconv.Itoa(read_type)}, "key": {key}}
		body := strings.NewReader(data.Encode())
		request_begin, err := http.NewRequest("POST", request_url, body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		request_begin = addHeader(request_begin)
		_, err = client.Do(request_begin)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// commit finish request
		// change the addtime as 90
		data.Set("addtime", strconv.Itoa(91))
		body = strings.NewReader(data.Encode())
		request_finish, err := http.NewRequest("POST", request_url, body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		request_finish = addHeader(request_finish)
		resp, err := client.Do(request_finish)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func Signin(key string) error {
	client := &http.Client{}
	data := url.Values{"route": {"signin"}, "key": {key}}
	body := strings.NewReader(data.Encode())
	request, err := http.NewRequest("POST", request_url, body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	request = addHeader(request)
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func PVP(key string) error {
	client := &http.Client{}
	data := url.Values{"route": {"get_counterpart"}, "key": {key}, "counter": {"0"}, "find_type": {"0"}}
	//add := 0
	var game_key = ""
	respData := make(map[string]interface{})
	question_bag := make(map[string]interface{})
	i := 0
	for {
		time.Sleep(1)
		data.Set("counter", strconv.Itoa(i))
		body := strings.NewReader(data.Encode())
		request, err := http.NewRequest("POST", request_url, body)
		if err != nil {
			return err
		}
		request = addHeader(request)
		response, _ := client.Do(request)
		defer response.Body.Close()
		resp_body, _ := io.ReadAll(response.Body)
		err = json.Unmarshal(resp_body, &respData)
		if err != nil {
			return err
		}
		if fmt.Sprint(respData["status"]) == "1" {
			question_bag = respData["question_bag"].(map[string]interface{})
			game_key = question_bag["gaming_key"].(string)
			break
		}
		i = i + 1
		if i > 10 {
			i = 0
		}
	}

	// Get the question Array
	questionArr := question_bag["question_arr"].([]interface{})
	// Get the question's number
	question_num := make([]string, 0)
	for _, question := range questionArr {
		question_num = append(question_num, question.(map[string]interface{})["num"].(string))
	}
	// Set Rand Seed
	rand.Seed(time.Now().UnixNano())

	// Get Answer of the question
	data_ask_question := url.Values{"route": {"ascertain_answer"}, "key": {key}, "gaming_key": {game_key}, "question_id": {".coctl"}, "answer_id": {""}, "question_num": {""}, "current_time": {"0"}}

	// send the ping for server
	go func() {
		ping_data := url.Values{"route": {"ask_opponent_score"}, "key": {key}, "gaming_key": {game_key}}
		ping_body := strings.NewReader(ping_data.Encode())
		request, _ := http.NewRequest("POST", request_url, ping_body)
		request = addHeader(request)
		for j := 0; j < 150; j++ {
			time.Sleep(1)
			response, _ := client.Do(request)
			jsonMap := make(map[string]interface{})
			response_body, _ := io.ReadAll(response.Body)
			err := json.Unmarshal(response_body, &jsonMap)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// range the question
	for index, num := range question_num {
		// Set the question_num for requestion body
		data_ask_question.Set("question_num", fmt.Sprint(num))
		body_ask_question := strings.NewReader(data_ask_question.Encode())
		request_ask_question, err := http.NewRequest("POST", request_url, body_ask_question)
		if err != nil {
			return err
		}
		request_ask_question = addHeader(request_ask_question)
		response_ask_question, err := client.Do(request_ask_question)
		defer response_ask_question.Body.Close()
		if err != nil {
			return err
		}

		// parse the response
		resp_body, err := io.ReadAll(response_ask_question.Body)
		err = json.Unmarshal(resp_body, &respData)
		if err != nil {
			return err
		}

		// Get the answer
		answer := respData["test_item2"].(map[string]interface{})["answer"]
		// Commit the answer
		data_ask_question.Set("answer_id", answer.(string))
		// Set the current time
		data_ask_question.Set("current_time", fmt.Sprint(rand.Intn(150)))
		body_commit_answer := strings.NewReader(data_ask_question.Encode())
		// commit the request
		request_commit_answer, err := http.NewRequest("POST", request_url, body_commit_answer)
		if err != nil {
			return err
		}
		request_commit_answer = addHeader(request_commit_answer)
		response_commit_answer, err := client.Do(request_commit_answer)
		if err != nil {
			return err
		}
		response_commit_answer.Body.Close()
		fmt.Println("完成了第", index, "题")
		time.Sleep(1)
	}

	// commit the end request
	time_out_data := url.Values{"route": {"time_out"}, "key": {key}, "gaming_key": {game_key}}
	time_out_body := strings.NewReader(time_out_data.Encode())
	request_time_out, err := http.NewRequest("POST", request_url, time_out_body)
	if err != nil {
		return err
	}
	request_time_out = addHeader(request_time_out)
	response_time_out, err := client.Do(request_time_out)
	if err != nil {
		return err
	}
	jsonMap := make(map[string]interface{})
	body, _ := io.ReadAll(response_time_out.Body)
	defer response_time_out.Body.Close()
	err = json.Unmarshal(body, &jsonMap)
	return nil
}

// Exec export the interface for user to exec the all function
func Exec(key string) (map[string]interface{}, error, error) {
	// get the user info first to check whether the key is valid
	err, info := GetUserInfo(key)
	if err != nil {
		return nil, err, errors.New("经典诵读的key已过期，请检查是否输入错误，若确认无误，请重新抓取key并发送给开发者")
	}
	fmt.Println(info)
	// Signin
	err = Signin(key)
	if err != nil {
		return nil, err, errors.New("登录错误，请联系开发者提交错误")
	}
	fmt.Println("登录成功")
	// per day question
	for i := 0; i < 3; i++ {
		err = Per_day_question(key)
		if err != nil {
			return nil, err, errors.New("执行到第 " + strconv.Itoa(i+1) + " 次每日一题错误，请联系开发者提交错误")
		}
		fmt.Println("第 " + strconv.Itoa(i+1) + " 次每日一题完成")
	}
	// read the poetry for 90sec
	err = Read(key)
	if err != nil {
		return nil, err, errors.New("阅读错误，请联系开发者提交错误")
	}
	fmt.Println("阅读完成")
	// pvp
	for i := 0; i < 3; i++ {
		err = PVP(key)
		if err != nil {
			return nil, err, errors.New("执行到第" + strconv.Itoa(i+1) + "次pvp出现问题，请检查key或者联系开发者")
		}
		fmt.Println("第" + strconv.Itoa(i+1) + "次PVP完成")
	}
	// get the info again
	err, info = GetUserInfo(key)
	if err != nil {
		return nil, err, errors.New("再次获得用户信息出错，请检查key或者联系开发者")
	}
	return info, nil, nil
}
