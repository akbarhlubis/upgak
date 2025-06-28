package notification

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier handles desktop notifications
type Notifier struct {
	Silent bool
}

// NewNotifier creates a new Notifier instance
func NewNotifier(silent bool) *Notifier {
	return &Notifier{Silent: silent}
}

// SendDownNotification sends a notification when a website goes down
func (n *Notifier) SendDownNotification(url, message string) error {
	if n.Silent {
		return nil
	}

	title := "Website Down"
	body := fmt.Sprintf("%s is down: %s", url, message)
	
	return n.sendNotification(title, body)
}

// SendUpNotification sends a notification when a website comes back up
func (n *Notifier) SendUpNotification(url string) error {
	if n.Silent {
		return nil
	}

	title := "Website Up"
	body := fmt.Sprintf("%s is back online", url)
	
	return n.sendNotification(title, body)
}

func (n *Notifier) sendNotification(title, body string) error {
	switch runtime.GOOS {
	case "linux":
		return n.sendLinuxNotification(title, body)
	case "darwin":
		return n.sendMacNotification(title, body)
	case "windows":
		return n.sendWindowsNotification(title, body)
	default:
		// Fallback: just print to console
		fmt.Printf("NOTIFICATION: %s - %s\n", title, body)
		return nil
	}
}

func (n *Notifier) sendLinuxNotification(title, body string) error {
	cmd := exec.Command("notify-send", title, body)
	return cmd.Run()
}

func (n *Notifier) sendMacNotification(title, body string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, body, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func (n *Notifier) sendWindowsNotification(title, body string) error {
	// For Windows, we'll use a PowerShell command
	script := fmt.Sprintf(`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null; $template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02); $toastXml = [xml] $template.GetXml(); $toastXml.GetElementsByTagName("text").AppendChild($toastXml.CreateTextNode("%s")) > $null; $toastXml.GetElementsByTagName("text").AppendChild($toastXml.CreateTextNode("%s")) > $null; [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("upgak").Show($template)`, title, body)
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}