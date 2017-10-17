
// FootballRobotDlg.h : ͷ�ļ�
//

#pragma once
#include "MyListCtrl.h"
#include "MyWebSocket.h"
// CFootballRobotDlg �Ի���
class CFootballRobotDlg : public CDialogEx
{
    // ����
public:
    CFootballRobotDlg ( CWnd* pParent = NULL ); // ��׼���캯��

    // �Ի�������
    enum { IDD = IDD_FOOTBALLROBOT_DIALOG };

protected:
    virtual void DoDataExchange ( CDataExchange* pDX ); // DDX/DDV ֧��


    // ʵ��
protected:
    HICON m_hIcon;

    // ���ɵ���Ϣӳ�亯��
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
