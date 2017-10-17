#pragma once
#include <string>
#include <vector>
#include <list>
#include "./libwebsockets/win32/include/libwebsockets.h"

enum State
{
    kStateConnecting = 0,   //连接中
    kStateOpen,             //已连接
    kStateClosing,          //关闭中
    kStateClosed            //已关闭
};

struct Data
{
    char* m_pBytes;     //串
    int m_nLen;         //长度
    bool m_bIsBinary;   //是否为二进制

    Data()
    {
        m_pBytes = NULL;
        m_nLen = 0;
        m_bIsBinary = false;
    }
};

class CMyWebSocket
{
public:
    class Delegate
    {
    public:
        virtual ~Delegate() {}
        virtual void onOpen ( CMyWebSocket* ws ) = 0;
        virtual void onMessage ( CMyWebSocket* ws, const Data& data ) = 0;
        virtual void onClose ( CMyWebSocket* ws ) = 0;
        virtual void onError ( CMyWebSocket* ws, const CMyWebSocket& error ) = 0;
    };

public:
    CMyWebSocket ( void );
    ~CMyWebSocket ( void );

public:
    State GetReadyState();
    bool Init();

private:
    State        m_readyState;
    std::string  m_sHost;
    unsigned int m_nPort;
    std::string  m_sPath;

    struct libwebsocket*         m_wsInstance;
    struct libwebsocket_context* m_wsContext;
    Delegate* _delegate;
    int m_nSSLConnection;
    struct libwebsocket_protocols* m_wsProtocols;
};

