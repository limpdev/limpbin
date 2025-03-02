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



# Transparent

Less than 1 minute

* * *

# Transparent

Create a transparent frameless window with custom chrome.
    
    
    import webview
    
    html = """
    <!doctype html>
    <html lang="en">
    	<head>
    		<meta charset="utf-8">
            <title>Test app</title>
            <style>
                .frame {
                    border-radius: 5px 5px 0 0;
                    position: fixed;
                    box-sizing: border-box;
                    width: 90%;
                    height: 90%;
                    background-color: #0055e4;
                    box-shadow: inset 1px 1px 1px 0px rgba(255,255,255,.25), inset -1px -1px 1px 0px rgba(0,0,0,.25), inset 0px 2px 4px -2px rgba(255,255,255,1);
                }
                .frame>tbody>tr>td {
                    vertical-align: top;
                }
                .header {
                    box-sizing: border-box;
                    padding: 5px;
                    height: 20px;
                    font-weight: bold;
                    color: white;
                }
                .header>img {
                    height: 16px;
                    transform: translateY(3px);
                }
                .content {
                    box-sizing: border-box;
                    background-color: #f0f0e8;
                    margin: 0 5px 5px 5px ;
                }
                .bodypanel {
                    background-color: #f0f0e8;
                    height: 100%;
                    box-shadow: 1px 1px 1px 0px rgba(255,255,255,.25), -1px -1px 1px 0px rgba(0,0,0,.25), inset 0px 0px 3px -2px rgba(0,0,0,1);
                    padding: 5px;
                }
            </style>
    	</head>
    	<body>
            <table class="frame">
                <tbody>
                    <tr>
                        <td class="header">
                            <img src="folder.png"/>
                            Danger!
                        </td>
                    </tr>
                    <tr>
                        <td class="body">
                            <div class="bodypanel">
                                <b>Alert!</b><br>
                                Lorem ipsum dolor sit amet, consectetur adipiscing elit
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
    	</body>
    </html>
    """
    
    if __name__ == '__main__':
        # Create a transparent webview window
        webview.create_window('Transparent window', html=html, transparent=True, frameless=True)
        webview.start()

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/transparent.md)

[PrevToggle Fullscreen](/examples/toggle_fullscreen)[NextUser Agent](/examples/user_agent)
