// MyListCtrl.cpp : 实现文件
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



// CMyListCtrl 消息处理程序

// 自适应列宽
void CMyListCtrl::AutoSize()
{
	// 禁用重绘,防止闪烁
	SetRedraw(FALSE);
	
	// 开始调整宽度
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

	// 开启重绘
	SetRedraw(TRUE);
}