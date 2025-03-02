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



# Dom Events

Less than 1 minute

* * *

# Dom Events

This example demonstrates how to expose Python functions to the Javascript domain.
    
    
    import webview
    from webview.dom import DOMEventHandler
    
    window = None
    
    def click_handler(e):
        print(e)
    
    def input_handler(e):
        print(e['target']['value'])
    
    def remove_handlers(scroll_event, click_event, input_event):
        scroll_event -= scroll_handler
        click_event -= click_handler
        input_event -= input_handler
    
    def scroll_handler(e):
        scroll_top = window.dom.window.node['scrollY']
        print(f'Scroll position {scroll_top}')
    
    def link_handler(e):
        print(f"Link target is {e['target']['href']}")
    
    def bind(window):
        window.dom.document.events.scroll += DOMEventHandler(scroll_handler, debounce=100)
    
        button = window.dom.get_element('#button')
        button.events.click += click_handler
    
        input = window.dom.get_element('#input')
        input.events.input += input_handler
    
        remove_events = window.dom.get_element('#remove')
        remove_events.on('click', lambda e: remove_handlers(window.dom.document.events.scroll, button.events.click, input.events.input))
    
        link = window.dom.get_element('#link')
        link.events.click += DOMEventHandler(link_handler, prevent_default=True)
    
    if __name__ == '__main__':
        window = webview.create_window(
            'DOM Event Example', html='''
                <html>
                    <head>
                    <style>
                        button {
                            font-size: 100%;
                            padding: 0.5rem;
                            margin: 0.3rem;
                            text-transform: uppercase;
                        }
                    </style>
                    </head>
                    <body style="height: 200vh;">
                        <div>
                            <input id="input" placeholder="Enter text">
                            <button id="button">Click me</button>
                            <a id="link" href="https://pywebview.flowrl.com">Click me</a>
                        </div>
                        <button id="remove" style="margin-top: 1rem;">Remove events</button>
                    </body>
                </html>
            '''
        )
        webview.start(bind, window)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/dom_events.md)

[PrevDestroy Window](/examples/destroy_window)[NextDom Manipulation](/examples/dom_manipulation)
