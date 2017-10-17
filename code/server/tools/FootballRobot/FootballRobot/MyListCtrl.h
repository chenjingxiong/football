#pragma once


// CMyListCtrl

class CMyListCtrl : public CListCtrl
{
	DECLARE_DYNAMIC(CMyListCtrl)

public:
	CMyListCtrl();
	virtual ~CMyListCtrl();

public:
	void AutoSize();	// 自适应列长度

protected:
	DECLARE_MESSAGE_MAP()
};


