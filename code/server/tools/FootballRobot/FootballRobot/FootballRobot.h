
// FootballRobot.h : PROJECT_NAME Ӧ�ó������ͷ�ļ�
//

#pragma once

#ifndef __AFXWIN_H__
	#error "�ڰ������ļ�֮ǰ������stdafx.h�������� PCH �ļ�"
#endif

#include "resource.h"		// ������


// CFootballRobotApp:
// �йش����ʵ�֣������ FootballRobot.cpp
//

class CFootballRobotApp : public CWinApp
{
public:
	CFootballRobotApp();

// ��д
public:
	virtual BOOL InitInstance();

// ʵ��

	DECLARE_MESSAGE_MAP()
};

extern CFootballRobotApp theApp;