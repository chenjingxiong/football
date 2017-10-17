// MyListCtrl.cpp : ʵ���ļ�
//

#include "stdafx.h"
#include "FootballRobot.h"
#include "MyListCtrl.h"


// CMyListCtrl

IMPLEMENT_DYNAMIC(CMyListCtrl, CLinkCtrl)

CMyListCtrl::CMyListCtrl()
{

}

CMyListCtrl::~CMyListCtrl()
{
}


BEGIN_MESSAGE_MAP(CMyListCtrl, CLinkCtrl)
END_MESSAGE_MAP()



// CMyListCtrl ��Ϣ�������

// ����Ӧ�п�
void CMyListCtrl::AutoSize()
{
	// �����ػ�,��ֹ��˸
	SetRedraw(FALSE);
	
	// ��ʼ�������
	CHeaderCtrl* pHeaderCtrl = GetHeaderCtrl();
	int nColNum = pHeaderCtrl->GetItemCount();;
	for (int i = 0; i < nColNum; i++)
	{
		SetColumnWidth(i, LVSCW_AUTOSIZE);
		int nColumnWidth = GetColumnWidth(i);
		SetColumnWidth(i, LVSCW_AUTOSIZE_USEHEADER);
		int nHeaderWidth = GetColumnWidth(i); 
		SetColumnWidth(i, max(nColumnWidth, nHeaderWidth) + 2);
	}

	// �����ػ�
	SetRedraw(TRUE);
}