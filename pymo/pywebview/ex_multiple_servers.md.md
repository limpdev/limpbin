Skip to main content

[![](..\\..\\..\\logo-no-text.png)pywebview](/)

[Guide](/guide/)

[API](/api/)

[Examples](/examples/)

[Contributing](/contributing/)

[Blog](/blog/)

[Changelog](/changelog)

[2.x](https://pywebview.flowrl.com/2.4)

[3.x](https://pywebview.flowrl.com/3.7)

[](https://github.com/r0x0r/pywebview)

  * [Cef](/examples/cef)
  * [Change Url](/examples/change_url)
  * [Confirm Close](/examples/confirm_close)
  * [Confirmation Dialog](/examples/confirmation_dialog)
  * [Cookies](/examples/cookies)
  * [Debug](/examples/debug)
  * [Destroy Window](/examples/destroy_window)
  * [Dom Events](/examples/dom_events)
  * [Dom Manipulation](/examples/dom_manipulation)
  * [Dom Traversal](/examples/dom_traversal)
  * [Downloads](/examples/downloads)
  * [Drag Drop](/examples/drag_drop)
  * [Drag Region](/examples/drag_region)
  * [Evaluate Js](/examples/evaluate_js)
  * [Evaluate Js Async](/examples/evaluate_js_async)
  * [Events](/examples/events)
  * [Expose](/examples/expose)
  * [Focus](/examples/focus)
  * [Frameless](/examples/frameless)
  * [Fullscreen](/examples/fullscreen)
  * [Get Current Url](/examples/get_current_url)
  * [Get Elements](/examples/get_elements)
  * [Hide Window](/examples/hide_window)
  * [Http Server](/examples/http_server)
  * [Icon](/examples/icon)
  * [Js Api](/examples/js_api)
  * [Links](/examples/links)
  * [Load Css](/examples/load_css)
  * [Load Html](/examples/load_html)
  * [Loading Animation](/examples/loading_animation)
  * [Localhost Ssl](/examples/localhost_ssl)
  * [Localization](/examples/localization)
  * [Menu](/examples/menu)
  * [Min Size](/examples/min_size)
  * [Move Window](/examples/move_window)
  * [Multiple Servers](/examples/multiple_servers)
  * [Multiple Windows](/examples/multiple_windows)
  * [On Top](/examples/on_top)
  * [Open File Dialog](/examples/open_file_dialog)
  * [Py2app Setup](/examples/py2app_setup)
  * [Pystray Icon](/examples/pystray_icon)
  * [Qt Test](/examples/qt_test)
  * [Remote Debugging](/examples/remote_debugging)
  * [Resize](/examples/resize)
  * [Run Js](/examples/run_js)
  * [Save File Dialog](/examples/save_file_dialog)
  * [Screens](/examples/screens)
  * [Settings](/examples/settings)
  * [Simple Browser](/examples/simple_browser)
  * [Toggle Fullscreen](/examples/toggle_fullscreen)
  * [Transparent](/examples/transparent)
  * [User Agent](/examples/user_agent)
  * [Vibrancy](/examples/vibrancy)
  * [Window State](/examples/window_state)
  * [Window Title Change](/examples/window_title_change)



# Multiple Servers

About 1 min

* * *

# Multiple Servers

Create multiple windows, some of which have their own servers, both before and after `webview.start()` is called.
    
    
    import bottle
    
    import webview
    
    # We'll have a global list of our windows so our web app can give us information
    # about them
    windows = []
    
    # A simple function to format a description of our servers
    def serverDescription(server):
        return f"{str(server).replace('<','').replace('>','')}"
    
    # Define a couple of simple web apps using Bottle
    app1 = bottle.Bottle()
    
    @app1.route('/')
    def hello():
        return '<h1>Second Window</h1><p>This one is a web app and has its own server.</p>'
    
    app2 = bottle.Bottle()
    
    @app2.route('/')
    def hello():
        head = """  <head>
                        <style type="text/css">
                            table {
                              font-family: arial, sans-serif;
                              border-collapse: collapse;
                              width: 100%;
                            }
    
                            td, th {
                              border: 1px solid #dddddd;
                              text-align: center;
                              padding: 8px;
                            }
    
                            tr:nth-child(even) {
                              background-color: #dddddd;
                            }
                        </style>
                    </head>
                """
        body = f""" <body>
                        <h1>Third Window</h1>
                        <p>This one is another web app and has its own server. It was started after webview.start.</p>
                        <p>Server Descriptions: </p>
                        <table>
                            <tr>
                                <th>Window</th>
                                <th>Object</th>
                                <th>IP Address</th>
                            </tr>
                            <tr>
                                <td>Global Server</td>
                                <td>{serverDescription(webview.http.global_server)}</td>
                                <td>{webview.http.global_server.address if not webview.http.global_server is None else 'None'}</td>
                            </tr>
                            <tr>
                                <td>First Window</td>
                                <td>{serverDescription(windows[0]._server)}</td>
                                <td>{windows[0]._server.address if not windows[0]._server is None else 'None'}</td>
                            </tr>
                            <tr>
                                <td>Second Window</td>
                                <td>{serverDescription(windows[1]._server)}</td>
                                <td>{windows[1]._server.address}</td>
                            </tr>
                            <tr>
                                <td>Third Window</td>
                                <td>{serverDescription(windows[2]._server)}</td>
                                <td>{windows[2]._server.address}</td>
                            </tr>
                        </table>
                    </body>
                """
        return head + body
    
    def third_window():
        # Create a new window after the loop started
        windows.append(webview.create_window('Window #3', url=app2))
    
    if __name__ == '__main__':
        # Master window
        windows.append(
            webview.create_window(
                'Window #1',
                html='<h1>First window</h1><p>This one is static HTML and just uses the global server for api calls.</p>',
            )
        )
        windows.append(webview.create_window('Window #2', url=app1, http_port=3333))
        webview.start(third_window, http_server=True, http_port=3334)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/multiple_servers.md)

[PrevMove Window](/examples/move_window)[NextMultiple Windows](/examples/multiple_windows)
