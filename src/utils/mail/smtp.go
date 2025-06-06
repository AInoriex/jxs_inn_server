package mail

import (
    "fmt"
    gomail "gopkg.in/mail.v2"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"math/rand"
)

// 生成随机6位验证码
func GenerateRandomEmailCode() string {
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 发送邮件
func SendEmail(addressee string, title string, text string) error {
	// 检查SMTP配置是否设置
    if config.CommonConfig.Smtp.Host == "" || config.CommonConfig.Smtp.Port <= 0 || config.CommonConfig.Smtp.Username == "" || config.CommonConfig.Smtp.Password == "" {
        log.Errorf("SendEmail SMTP配置字段有误")
		return fmt.Errorf("SMTP configuration is not set up correctly")
    }
    // 检查参数
    if addressee == "" || title == "" || text == ""{
        log.Errorf("SendEmail 参数有误, addressee: %s, title: %s, text: %s", addressee, title, text)
		return fmt.Errorf("SendEmail parameters error")
    }

    // Create a new message
    message := gomail.NewMessage()
    message.SetHeader("From", config.CommonConfig.Smtp.Username) // 发件人邮箱
    message.SetHeader("To", addressee) // 收件人邮箱
    message.SetHeader("Subject", "Hello from the AInoriex") // 邮件主题
    message.SetBody("text/plain", "Hello World.") // 邮件内容

    // Set up the SMTP dialer
    dialer := gomail.NewDialer(config.CommonConfig.Smtp.Host, config.CommonConfig.Smtp.Port, config.CommonConfig.Smtp.Username, config.CommonConfig.Smtp.Password)

    // Send the email
    if err := dialer.DialAndSend(message); err != nil {
        log.Errorf("SendEmail failed to send email: %v", err)
        return err
    } 
	log.Infof("SendEmail sent successfully")
	return nil
}