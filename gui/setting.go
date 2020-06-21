package gui

import (
	"github.com/Asutorufa/yuhaiin/config"
	"github.com/Asutorufa/yuhaiin/process"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type setting struct {
	settingWindow *widgets.QMainWindow
	parent        *widgets.QMainWindow

	BlackIconCheckBox         *widgets.QCheckBox
	DnsOverHttpsCheckBox      *widgets.QCheckBox
	bypassCheckBox            *widgets.QCheckBox
	DnsOverHttpsProxyCheckBox *widgets.QCheckBox

	redirProxyAddressLabel   *widgets.QLabel
	httpAddressLabel         *widgets.QLabel
	socks5BypassAddressLabel *widgets.QLabel
	dnsServerLabel           *widgets.QLabel
	ssrPathLabel             *widgets.QLabel
	//BypassFileLabel          *widgets.QLabel
	dnsSubNetLabel *widgets.QLabel

	redirProxyAddressLineText *widgets.QLineEdit
	httpAddressLineText       *widgets.QLineEdit
	socks5BypassLineText      *widgets.QLineEdit
	dnsServerLineText         *widgets.QLineEdit
	ssrPathLineText           *widgets.QLineEdit
	BypassFileLineText        *widgets.QLineEdit
	dnsSubNetLineText         *widgets.QLineEdit

	applyButton      *widgets.QPushButton
	updateRuleButton *widgets.QPushButton
}

func NewSettingWindow(parent *widgets.QMainWindow) *widgets.QMainWindow {
	s := setting{}
	s.parent = parent
	s.settingWindow = widgets.NewQMainWindow(nil, core.Qt__Window)
	s.settingWindow.SetWindowTitle("setting")
	s.settingWindow.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		event.Ignore()
		s.settingWindow.Hide()
	})
	s.settingInit()
	s.setLayout()
	//s.setGeometry()
	s.setListener()
	s.extends()

	return s.settingWindow
}

func (s *setting) settingInit() {
	s.BlackIconCheckBox = widgets.NewQCheckBox2("Black Icon", s.settingWindow)
	s.DnsOverHttpsCheckBox = widgets.NewQCheckBox2("DOH", s.settingWindow)
	//httpProxyCheckBox := widgets.NewQCheckBox2("http proxy", sGui.settingWindow)
	s.bypassCheckBox = widgets.NewQCheckBox2("BYPASS", s.settingWindow)
	s.DnsOverHttpsProxyCheckBox = widgets.NewQCheckBox2("PROXY", s.settingWindow)
	s.redirProxyAddressLabel = widgets.NewQLabel2("REDIR", s.settingWindow, 0)
	s.redirProxyAddressLineText = widgets.NewQLineEdit(s.settingWindow)
	s.httpAddressLabel = widgets.NewQLabel2("HTTP", s.settingWindow, 0)
	s.httpAddressLineText = widgets.NewQLineEdit(s.settingWindow)
	s.socks5BypassAddressLabel = widgets.NewQLabel2("SOCKS5", s.settingWindow, 0)
	s.socks5BypassLineText = widgets.NewQLineEdit(s.settingWindow)
	s.dnsServerLabel = widgets.NewQLabel2("DNS", s.settingWindow, 0)
	s.dnsServerLineText = widgets.NewQLineEdit(s.settingWindow)
	s.ssrPathLabel = widgets.NewQLabel2("SSR PATH", s.settingWindow, 0)
	s.ssrPathLineText = widgets.NewQLineEdit(s.settingWindow)
	//s.BypassFileLabel = widgets.NewQLabel2("bypassFile", s.settingWindow, 0)
	s.BypassFileLineText = widgets.NewQLineEdit(s.settingWindow)
	s.applyButton = widgets.NewQPushButton2("apply", s.settingWindow)
	s.updateRuleButton = widgets.NewQPushButton2("Reimport Bypass Rule", s.settingWindow)
	s.dnsSubNetLabel = widgets.NewQLabel2("SUBNET", nil, 0)
	s.dnsSubNetLineText = widgets.NewQLineEdit(nil)
}

func (s *setting) setLayout() {
	localProxyGroup := widgets.NewQGroupBox2("PROXY", nil)
	localProxyLayout := widgets.NewQGridLayout2()
	localProxyLayout.AddWidget2(s.httpAddressLabel, 0, 0, 0)
	localProxyLayout.AddWidget2(s.httpAddressLineText, 0, 1, 0)
	localProxyLayout.AddWidget2(s.socks5BypassAddressLabel, 1, 0, 0)
	localProxyLayout.AddWidget2(s.socks5BypassLineText, 1, 1, 0)
	localProxyLayout.AddWidget2(s.redirProxyAddressLabel, 2, 0, 0)
	localProxyLayout.AddWidget2(s.redirProxyAddressLineText, 2, 1, 0)
	localProxyGroup.SetLayout(localProxyLayout)

	dnsGroup := widgets.NewQGroupBox2("DNS", nil)
	dnsLayout := widgets.NewQGridLayout2()
	dnsLayout.AddWidget2(s.DnsOverHttpsCheckBox, 0, 0, 0)
	dnsLayout.AddWidget2(s.DnsOverHttpsProxyCheckBox, 0, 1, 0)
	dnsLayout.AddWidget2(s.dnsServerLabel, 1, 0, 0)
	dnsLayout.AddWidget2(s.dnsServerLineText, 1, 1, 0)
	dnsLayout.AddWidget2(s.dnsSubNetLabel, 2, 0, 0)
	dnsLayout.AddWidget2(s.dnsSubNetLineText, 2, 1, 0)
	dnsGroup.SetLayout(dnsLayout)

	bypassGroup := widgets.NewQGroupBox2("BYPASS", nil)
	bypassLayout := widgets.NewQGridLayout2()
	bypassLayout.AddWidget2(s.bypassCheckBox, 0, 0, 0)
	bypassLayout.AddWidget2(s.BypassFileLineText, 1, 0, 0)
	bypassGroup.SetLayout(bypassLayout)

	othersGroup := widgets.NewQGroupBox2("OTHERS", nil)
	othersLayout := widgets.NewQGridLayout2()
	othersLayout.AddWidget3(s.BlackIconCheckBox, 0, 0, 1, 2, 0)
	othersLayout.AddWidget2(s.ssrPathLabel, 1, 0, 0)
	othersLayout.AddWidget2(s.ssrPathLineText, 1, 1, 0)
	othersGroup.SetLayout(othersLayout)

	windowLayout := widgets.NewQGridLayout2()
	windowLayout.AddWidget2(localProxyGroup, 0, 0, 0)
	windowLayout.AddWidget2(dnsGroup, 0, 1, 0)
	windowLayout.AddWidget2(bypassGroup, 1, 0, 0)
	windowLayout.AddWidget2(othersGroup, 1, 1, 0)
	windowLayout.AddWidget2(s.applyButton, 2, 0, 0)
	windowLayout.AddWidget2(s.updateRuleButton, 3, 0, 0)

	centralWidget := widgets.NewQWidget(s.settingWindow, 0)
	centralWidget.SetLayout(windowLayout)
	s.settingWindow.SetCentralWidget(centralWidget)
}

func (s *setting) setGeometry() {
	s.BlackIconCheckBox.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 0), core.NewQPoint2(140, 30)))
	//httpProxyCheckBox.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 40), core.NewQPoint2(130, 70)))
	s.redirProxyAddressLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 80), core.NewQPoint2(70, 110)))
	s.redirProxyAddressLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(80, 80), core.NewQPoint2(210, 110)))
	s.httpAddressLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 120), core.NewQPoint2(70, 150)))
	s.httpAddressLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(80, 120), core.NewQPoint2(210, 150)))
	s.socks5BypassAddressLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(220, 120), core.NewQPoint2(290, 150)))
	s.socks5BypassLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(300, 120), core.NewQPoint2(420, 150)))
	s.dnsServerLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 160), core.NewQPoint2(50, 190)))
	s.DnsOverHttpsCheckBox.SetGeometry(core.NewQRect2(core.NewQPoint2(60, 160), core.NewQPoint2(125, 190)))
	s.DnsOverHttpsProxyCheckBox.SetGeometry(core.NewQRect2(core.NewQPoint2(135, 160), core.NewQPoint2(220, 190)))
	s.dnsServerLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(230, 160), core.NewQPoint2(420, 190)))
	s.ssrPathLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 200), core.NewQPoint2(90, 230)))
	s.ssrPathLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(100, 200), core.NewQPoint2(420, 230)))
	s.bypassCheckBox.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 240), core.NewQPoint2(105, 270)))
	//s.BypassFileLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 240), core.NewQPoint2(100, 270)))
	s.BypassFileLineText.SetGeometry(core.NewQRect2(core.NewQPoint2(115, 240), core.NewQPoint2(420, 270)))
	s.applyButton.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 280), core.NewQPoint2(105, 310)))
	s.updateRuleButton.SetGeometry(core.NewQRect2(core.NewQPoint2(115, 280), core.NewQPoint2(300, 310)))
}

func (s *setting) setListener() {
	// Listen
	update := func() {
		conFig, err := config.SettingDecodeJSON()
		if err != nil {
			MessageBox(err.Error())
			return
		}
		s.BlackIconCheckBox.SetChecked(conFig.BlackIcon)
		s.DnsOverHttpsCheckBox.SetChecked(conFig.IsDNSOverHTTPS)
		//httpProxyCheckBox.SetChecked(conFig.HttpProxy)
		s.bypassCheckBox.SetChecked(conFig.Bypass)
		s.DnsOverHttpsProxyCheckBox.SetChecked(conFig.DNSAcrossProxy)
		s.redirProxyAddressLineText.SetText(conFig.RedirProxyAddress)
		s.httpAddressLineText.SetText(conFig.HttpProxyAddress)
		s.socks5BypassLineText.SetText(conFig.Socks5ProxyAddress)
		s.dnsServerLineText.SetText(conFig.DnsServer)
		s.ssrPathLineText.SetText(conFig.SsrPath)
		s.BypassFileLineText.SetText(conFig.BypassFile)
		s.dnsSubNetLineText.SetText(conFig.DnsSubNet)
	}

	applyClick := func(bool2 bool) {
		conFig, err := config.SettingDecodeJSON()
		if err != nil {
			MessageBox(err.Error())
			return
		}

		conFig.BlackIcon = s.BlackIconCheckBox.IsChecked()
		//conFig.HttpProxy = httpProxyCheckBox.IsChecked()
		conFig.Bypass = s.bypassCheckBox.IsChecked()
		conFig.IsDNSOverHTTPS = s.DnsOverHttpsCheckBox.IsChecked()
		conFig.DNSAcrossProxy = s.DnsOverHttpsProxyCheckBox.IsChecked()
		conFig.DnsServer = s.dnsServerLineText.Text()
		conFig.SsrPath = s.ssrPathLineText.Text()
		conFig.HttpProxyAddress = s.httpAddressLineText.Text()
		conFig.Socks5ProxyAddress = s.socks5BypassLineText.Text()
		conFig.RedirProxyAddress = s.redirProxyAddressLineText.Text()
		conFig.BypassFile = s.BypassFileLineText.Text()

		go func() {
			process.SetConFig(conFig)
		}()

		if err := config.SettingEnCodeJSON(conFig); err != nil {
			MessageBox(err.Error())
		}

		update()
		MessageBox("Applied.")
	}

	// set Listener
	s.applyButton.ConnectClicked(applyClick)
	s.updateRuleButton.ConnectClicked(func(checked bool) {
		if err := process.UpdateMatch(); err != nil {
			MessageBox(err.Error())
			return
		}
		MessageBox("Updated.")
	})

	s.settingWindow.ConnectShowEvent(func(event *gui.QShowEvent) {
		update()
	})
}
