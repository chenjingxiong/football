#pragma once


// CMyListCtrl

class CMyListCtrl : public CListCtrl
{
	DECLARE_DYNAMIC(CMyListCtrl)

public:
	CMyListCtrl();
	virtual ~CMyListCtrl();

public:
	void AutoSize();	// ����Ӧ�г���

protected:
	DECLARE_MESSAGE_MAP()
};


