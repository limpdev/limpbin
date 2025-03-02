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



# Dom Traversal

Less than 1 minute

* * *

# Dom Traversal

This example demonstrates how to traverse DOM in Python.
    
    
    import webview
    
    def bind(window):
        container = window.dom.get_element('#container')
        container_button = window.dom.get_element('#container-button')
        blue_rectangle = window.dom.get_element('#blue-rectangle')
        blue_parent_button = window.dom.get_element('#blue-parent-button')
        blue_next_button = window.dom.get_element('#blue-next-button')
        blue_previous_button = window.dom.get_element('#blue-previous-button')
    
        container_button.events.click += lambda e: print(container.children)
        blue_parent_button.events.click += lambda e: print(blue_rectangle.parent)
        blue_next_button.events.click += lambda e: print(blue_rectangle.next)
        blue_previous_button.events.click += lambda e: print(blue_rectangle.previous)
    
    if __name__ == '__main__':
        window = webview.create_window(
            'DOM Manipulations Example', html='''
                <html>
                    <head>
                    <style>
                        button {
                            font-size: 100%;
                            padding: 0.5rem;
                            margin: 0.3rem;
                            text-transform: uppercase;
                        }
    
                        .rectangle {
                            width: 100px;
                            height: 100px;
                            display: flex;
                            justify-content: center;
                            align-items: center;
                            color: white;
                            margin-right: 5px;
                        }
                    </style>
                    </head>
                    <body>
                        <h1>Container</h1>
                        <div id="container" style="border: 1px #eee solid; display: flex; padding: 10px 0;">
                            <div id="red-rectangle" class="rectangle" style="background-color: red;">RED</div>
                            <div id="blue-rectangle" class="rectangle" style="background-color: blue;">BLUE</div>
                            <div id="green-rectangle" class="rectangle" style="background-color: green;">GREEN</div>
                        </div>
                        <button id="container-button">Get container's children</button>
                        <button id="blue-parent-button">Get blue's parent</button>
                        <button id="blue-next-button">Get blue's next element</button>
                        <button id="blue-previous-button">Get blue's previous element</button>
                    </body>
                </html>
            '''
        )
        webview.start(bind, window)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/dom_traversal.md)

[PrevDom Manipulation](/examples/dom_manipulation)[NextDownloads](/examples/downloads)
