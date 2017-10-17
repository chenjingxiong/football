#include "stdafx.h"
#include "MyWebSocket.h"

CMyWebSocket::CMyWebSocket()
{
    // 初始化
    m_readyState = kStateConnecting;
    m_nPort = 8080;
    m_sHost = "192.168.20.54";
    m_sPath = "";
    m_nSSLConnection = 0;
    m_wsInstance = NULL;
    m_wsContext = NULL;
    m_wsProtocols = NULL;
}

CMyWebSocket::~CMyWebSocket()
{
}

static int
Test ( struct libwebsocket_context* context,
       struct libwebsocket* wsi,
       enum libwebsocket_callback_reasons reason,
       void* user, void* in, size_t len )
{
    switch ( reason )
    {
        case LWS_CALLBACK_DEL_POLL_FD:
        case LWS_CALLBACK_PROTOCOL_DESTROY:
        case LWS_CALLBACK_CLIENT_CONNECTION_ERROR:
            {
                //                 AfxMessageBox ( "LWS_CALLBACK_DEL_POLL_FD | \
                //                                 LWS_CALLBACK_PROTOCOL_DESTROY | LWS_CALLBACK_CLIENT_CONNECTION_ERROR" );
            }
            break;

        case LWS_CALLBACK_CLIENT_ESTABLISHED:
            {
                libwebsocket_callback_on_writable ( context, wsi );
            }
            break;

        case LWS_CALLBACK_CLIENT_RECEIVE:
            {
                AfxMessageBox ( "LWS_CALLBACK_CLIENT_RECEIVE" );
            }
            break;

        case LWS_CALLBACK_CLIENT_WRITEABLE:
            {
                int bytesWrite = 0;
                unsigned char buf[255] = "HelloWorld";
                bytesWrite = libwebsocket_write ( wsi,  buf, 255, LWS_WRITE_TEXT );

                if ( bytesWrite < 0 )
                {
                    AfxMessageBox ( "Wrong" );
                }

                libwebsocket_callback_on_writable ( context, wsi );
            }
            break;

        case LWS_CALLBACK_CLOSED:
            {
                AfxMessageBox ( "LWS_CALLBACK_CLOSED" );
            }
            break;

        case LWS_CALLBACK_ADD_POLL_FD:
            {
                //AfxMessageBox ( "LWS_CALLBACK_ADD_POLL_FD" );
            }
            break;

        case LWS_CALLBACK_CLIENT_CONFIRM_EXTENSION_SUPPORTED:

        //             if ( ( strcmp ( ( char* ) in, "deflate-stream" ) == 0 ) && deny_deflate ) {
        //                 fprintf ( stderr, "denied deflate-stream extension\n" );
        //                 return 1;
        //             }
        //
        //             if ( ( strcmp ( ( char* ) in, "x-google-mux" ) == 0 ) && deny_mux ) {
        //                 fprintf ( stderr, "denied x-google-mux extension\n" );
        //                 return 1;
        //             }
        default:
            break;
    }

    return 0;
}
static struct libwebsocket_protocols protocols[] = {
    {
        "dumb-increment-protocol",
        Test,
        0,
        255,
    },
    { NULL, NULL, 0, 0 } /* end */
};
// 初始化
bool CMyWebSocket::Init()
{
    struct lws_context_creation_info info;
    memset ( &info, 0, sizeof info );
    info.port = CONTEXT_PORT_NO_LISTEN;
    info.protocols = protocols;
    #ifndef LWS_NO_EXTENSIONS
    info.extensions = libwebsocket_get_internal_extensions();
    #endif
    info.gid = -1;
    info.uid = -1;
    info.user = ( void* ) this;
    struct libwebsocket_context* _wsContext;
    _wsContext = libwebsocket_create_context ( &info );

    if ( _wsContext == NULL ) {
        fprintf ( stderr, "Creating libwebsocket context failed\n" );
        return 1;
    }

    struct libwebsocket* _wsInstance;

    _wsInstance = libwebsocket_client_connect ( _wsContext, m_sHost.c_str(), m_nPort, m_nSSLConnection,
                  m_sPath.c_str(), m_sHost.c_str(), m_sHost.c_str(),
                  protocols->name, -1 );

    if ( _wsInstance == NULL ) {
        fprintf ( stderr, "connect failed\n" );
    }

    libwebsocket_service ( _wsContext, 0 );
    int bytesWrite = 0;
    unsigned char buf[255] = "HelloWorld";
    bytesWrite = libwebsocket_write ( _wsInstance,  buf, 12, LWS_WRITE_TEXT );

    if ( bytesWrite < 0 )
    {
        AfxMessageBox ( "Wrong" );
    }

    //     if ( _wsContext != NULL )
    //     {
    //         libwebsocket_context_destroy ( _wsContext );
    //     }
    return true;
}