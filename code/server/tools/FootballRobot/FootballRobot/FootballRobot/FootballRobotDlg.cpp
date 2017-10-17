
// FootballRobotDlg.cpp : ʵ���ļ�
//

#include "stdafx.h"
#include "FootballRobot.h"
#include "FootballRobotDlg.h"
#include "afxdialogex.h"

#ifdef _DEBUG
#define new DEBUG_NEW
#endif


// ����Ӧ�ó��򡰹��ڡ��˵���� CAboutDlg �Ի���

class CAboutDlg : public CDialogEx {
public:
    CAboutDlg();

    // �Ի�������
    enum { IDD = IDD_ABOUTBOX };

protected:
    virtual void DoDataExchange ( CDataExchange* pDX ); // DDX/DDV ֧��

    // ʵ��
protected:
    DECLARE_MESSAGE_MAP()
};

CAboutDlg::CAboutDlg() : CDialogEx ( CAboutDlg::IDD )
{
}

void CAboutDlg::DoDataExchange ( CDataExchange* pDX )
{
    CDialogEx::DoDataExchange ( pDX );
}

BEGIN_MESSAGE_MAP ( CAboutDlg, CDialogEx )
END_MESSAGE_MAP()


// CFootballRobotDlg �Ի���



CFootballRobotDlg::CFootballRobotDlg ( CWnd* pParent /*=NULL*/ )
    : CDialogEx ( CFootballRobotDlg::IDD, pParent )
    , m_edt_accountName ( _T ( "robot" ) )
    , m_edt_accountPwd ( _T ( "123456" ) )
    , m_edt_loginNum ( 1 )
    , m_edt_loginServer ( _T ( "192.168.20.54:8080" ) )
    , m_edt_actionDelay ( 500 )
{
    m_hIcon = AfxGetApp()->LoadIcon ( IDR_MAINFRAME );
}

void CFootballRobotDlg::DoDataExchange ( CDataExchange* pDX )
{
    CDialogEx::DoDataExchange ( pDX );
    DDX_Control ( pDX, IDL_ROBOTLIST, m_listCtrl );
    DDX_Text ( pDX, IDE_ROBOTACCOUNTNAME, m_edt_accountName );
    DDV_MaxChars ( pDX, m_edt_accountName, 255 );
    DDX_Text ( pDX, IDE_ROBOTACCOUNTPWD, m_edt_accountPwd );
    DDV_MaxChars ( pDX, m_edt_accountPwd, 255 );
    DDX_Text ( pDX, IDE_NUMBER, m_edt_loginNum );
    DDV_MinMaxInt ( pDX, m_edt_loginNum, 1, 9999 );
    DDX_Text ( pDX, IDE_LOGINSERVER, m_edt_loginServer );
    DDV_MaxChars ( pDX, m_edt_loginServer, 255 );
    DDX_Text ( pDX, IDE_ACTIONDELAY, m_edt_actionDelay );
    DDV_MinMaxInt ( pDX, m_edt_actionDelay, 1, 10000 );
}

BEGIN_MESSAGE_MAP ( CFootballRobotDlg, CDialogEx )
    ON_WM_SYSCOMMAND()
    ON_WM_PAINT()
    ON_WM_QUERYDRAGICON()
    ON_WM_SIZE()
END_MESSAGE_MAP()


// CFootballRobotDlg ��Ϣ��������

BOOL CFootballRobotDlg::OnInitDialog()
{
    CDialogEx::OnInitDialog();
    // ��������...���˵������ӵ�ϵͳ�˵��С�
    // IDM_ABOUTBOX ������ϵͳ���Χ�ڡ�
    ASSERT ( ( IDM_ABOUTBOX & 0xFFF0 ) == IDM_ABOUTBOX );
    ASSERT ( IDM_ABOUTBOX < 0xF000 );
    CMenu* pSysMenu = GetSystemMenu ( FALSE );

    if ( pSysMenu != NULL ) {
        BOOL bNameValid;
        CString strAboutMenu;
        bNameValid = strAboutMenu.LoadString ( IDS_ABOUTBOX );
        ASSERT ( bNameValid );

        if ( !strAboutMenu.IsEmpty() ) {
            pSysMenu->AppendMenu ( MF_SEPARATOR );
            pSysMenu->AppendMenu ( MF_STRING, IDM_ABOUTBOX, strAboutMenu );
        }
    }

    // ���ô˶Ի����ͼ�ꡣ��Ӧ�ó��������ڲ��ǶԻ���ʱ����ܽ��Զ�
    //  ִ�д˲���
    SetIcon ( m_hIcon, TRUE );      // ���ô�ͼ��
    SetIcon ( m_hIcon, FALSE );     // ����Сͼ��
    ShowWindow ( SW_NORMAL );
    // ��ʼ���б�ͷ��Ϣ
    m_listCtrl.InsertColumn ( 1, "ID" );        // ID
    m_listCtrl.InsertColumn ( 2, "AccountName" ); // �˺���
    m_listCtrl.InsertColumn ( 3, "AccountPwd" ); // �ʺ�����
    m_listCtrl.InsertColumn ( 4, "TeamID" );    // ����ID
    m_listCtrl.InsertColumn ( 5, "TeamName" );  // ������
    m_listCtrl.InsertColumn ( 6, "Coin" );      // ����
    m_listCtrl.InsertColumn ( 7, "Diamond" );   // ��ʯ
    m_listCtrl.InsertColumn ( 8, "Action" );    // ��Ϊ
    m_listCtrl.InsertColumn ( 9, "Upload" );    // �ϴ�
    m_listCtrl.InsertColumn ( 10, "Download" ); // ����
    m_listCtrl.AutoSize();
    CMyWebSocket web;
    web.Init();
    return TRUE;  // ���ǽ��������õ��ؼ������򷵻� TRUE
}

void CFootballRobotDlg::OnSysCommand ( UINT nID, LPARAM lParam )
{
    if ( ( nID & 0xFFF0 ) == IDM_ABOUTBOX ) {
        CAboutDlg dlgAbout;
        dlgAbout.DoModal();
    }
    else {
        CDialogEx::OnSysCommand ( nID, lParam );
    }
}

// �����Ի���������С����ť������Ҫ����Ĵ���
//  �����Ƹ�ͼ�ꡣ����ʹ���ĵ�/��ͼģ�͵� MFC Ӧ�ó���
//  �⽫�ɿ���Զ���ɡ�

void CFootballRobotDlg::OnPaint()
{
    if ( IsIconic() ) {
        CPaintDC dc ( this ); // ���ڻ��Ƶ��豸������
        SendMessage ( WM_ICONERASEBKGND, reinterpret_cast<WPARAM> ( dc.GetSafeHdc() ), 0 );
        // ʹͼ���ڹ����������о���
        int cxIcon = GetSystemMetrics ( SM_CXICON );
        int cyIcon = GetSystemMetrics ( SM_CYICON );
        CRect rect;
        GetClientRect ( &rect );
        int x = ( rect.Width() - cxIcon + 1 ) / 2;
        int y = ( rect.Height() - cyIcon + 1 ) / 2;
        // ����ͼ��
        dc.DrawIcon ( x, y, m_hIcon );
    }
    else {
        CDialogEx::OnPaint();
    }
}

//���û��϶���С������ʱϵͳ���ô˺���ȡ�ù��
//��ʾ��
HCURSOR CFootballRobotDlg::OnQueryDragIcon()
{
    return static_cast<HCURSOR> ( m_hIcon );
}