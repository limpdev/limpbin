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



# Menu

April 23, 2022Less than 1 minute

* * *

# Menu

Create an application menu.
    
    
    import webview
    import webview.menu as wm
    
    def change_active_window_content():
        active_window = webview.active_window()
        if active_window:
            active_window.load_html('<h1>You changed this window!</h1>')
    
    def click_me():
        active_window = webview.active_window()
        if active_window:
            active_window.load_html('<h1>You clicked me!</h1>')
    
    def do_nothing():
        pass
    
    def say_this_is_window_2():
        active_window = webview.active_window()
        if active_window:
            active_window.load_html('<h1>This is window 2</h2>')
    
    def open_file_dialog():
        active_window = webview.active_window()
        active_window.create_file_dialog(webview.SAVE_DIALOG, directory='/', save_filename='test.file')
    
    if __name__ == '__main__':
        window_1 = webview.create_window(
            'Application Menu Example', 'https://pywebview.flowrl.com/hello'
        )
        window_2 = webview.create_window(
            'Another Window', html='<h1>Another window to test application menu</h1>'
        )
    
        menu_items = [
            wm.Menu(
                'Test Menu',
                [
                    wm.MenuAction('Change Active Window Content', change_active_window_content),
                    wm.MenuSeparator(),
                    wm.Menu(
                        'Random',
                        [
                            wm.MenuAction('Click Me', click_me),
                            wm.MenuAction('File Dialog', open_file_dialog),
                        ],
                    ),
                ],
            ),
            wm.Menu('Nothing Here', [wm.MenuAction('This will do nothing', do_nothing)]),
        ]
    
        webview.start(menu=menu_items)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/menu.md)

Last update: 9/26/2023, 6:21:55 AM

[PrevLocalization](/examples/localization)[NextMin Size](/examples/min_size)
