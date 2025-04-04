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



# Events

July 7, 2019Less than 1 minute

* * *

# Events

Subscribe and unsubscribe to pywebview events.
    
    
    import webview
    
    def on_before_show(window):
        print('Native window object', window.native)
    
    def on_closed():
        print('pywebview window is closed')
    
    def on_closing():
        print('pywebview window is closing')
    
    def on_shown():
        print('pywebview window shown')
    
    def on_minimized():
        print('pywebview window minimized')
    
    def on_restored():
        print('pywebview window restored')
    
    def on_maximized():
        print('pywebview window maximized')
    
    def on_resized(width, height):
        print(
            'pywebview window is resized. new dimensions are {width} x {height}'.format(
                width=width, height=height
            )
        )
    
    # you can supply optional window argument to access the window object event was triggered on
    def on_loaded(window):
        print('DOM is ready')
    
        # unsubscribe event listener
        window.events.loaded -= on_loaded
        window.load_url('https://pywebview.flowrl.com/hello')
    
    def on_moved(x, y):
        print('pywebview window is moved. new coordinates are x: {x}, y: {y}'.format(x=x, y=y))
    
    if __name__ == '__main__':
        window = webview.create_window(
            'Simple browser', 'https://pywebview.flowrl.com/', confirm_close=True
        )
    
        window.events.closed += on_closed
        window.events.closing += on_closing
        window.events.before_show += on_before_show
        window.events.shown += on_shown
        window.events.loaded += on_loaded
        window.events.minimized += on_minimized
        window.events.maximized += on_maximized
        window.events.restored += on_restored
        window.events.resized += on_resized
        window.events.moved += on_moved
    
        webview.start()

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/events.md)

Last update: 9/26/2023, 6:21:55 AM

[PrevEvaluate Js Async](/examples/evaluate_js_async)[NextExpose](/examples/expose)
