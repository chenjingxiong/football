
// FootballRobotDlg.h : 头文件
//

#pragma once
#include "MyListCtrl.h"
#include "MyWebSocket.h"
// CFootballRobotDlg 对话框
class CFootballRobotDlg : public CDialogEx
{
    // 构造
public:
    CFootballRobotDlg ( CWnd* pParent = NULL ); // 标准构造函数

    // 对话框数据
    enum { IDD = IDD_FOOTBALLROBOT_DIALOG };

protected:
    virtual void DoDataExchange ( CDataExchange* pDX ); // DDX/DDV 支持


    // 实现
protected:
    HICON m_hIcon;

    // 生成的消息映射函数
    virtual BOOL OnInitDialog();
    afx_msg void OnSysCommand ( UINT nID, LPARAM lParam );
    afx_msg void OnPaint();
    afx_msg HCURSOR OnQueryDragIcon();
    DECLARE_MESSAGE_MAP()

public:
    CMyListCtrl m_listCtrl;
    CString m_edt_accountName;
    CString m_edt_accountPwd;
    int m_edt_loginNum;
    CString m_edt_loginServer;
    int m_edt_actionDelay;
};
