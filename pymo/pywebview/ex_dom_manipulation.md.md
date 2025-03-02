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



# Dom Manipulation

About 1 min

* * *

# Dom Manipulation

This example demonstrates how to manipulate DOM in Python.
    
    
    import random
    import webview
    
    rectangles = []
    
    def random_color():
        red = random.randint(0, 255)
        green = random.randint(0, 255)
        blue = random.randint(0, 255)
    
        return f'rgb({red}, {green}, {blue})'
    
    def bind(window):
        def toggle_disabled():
            disabled = None if len(rectangles) > 0 else True
            remove_button.attributes = { 'disabled': disabled }
            empty_button.attributes = { 'disabled': disabled }
            move_button.attributes = { 'disabled': disabled }
    
        def create_rectangle(_):
            color = random_color()
            rectangle = window.dom.create_element(f'<div class="rectangle" style="background-color: {color};"></div>', rectangle_container)
            rectangles.append(rectangle)
            toggle_disabled()
    
        def remove_rectangle(_):
            if len(rectangles) > 0:
                rectangles.pop().remove()
            toggle_disabled()
    
        def move_rectangle(_):
            if len(rectangle_container.children) > 0:
                rectangle_container.children[-1].move(circle_container)
    
        def empty_container(_):
            rectangle_container.empty()
            rectangles.clear()
            toggle_disabled()
    
        def change_color(_):
            circle.style['background-color'] = random_color()
    
        def toggle_class(_):
            circle.classes.toggle('circle')
    
        rectangle_container = window.dom.get_element('#rectangles')
        circle_container = window.dom.get_element('#circles')
        circle = window.dom.get_element('#circle')
    
        toggle_button = window.dom.get_element('#toggle-button')
        toggle_class_button = window.dom.get_element('#toggle-class-button')
        duplicate_button = window.dom.get_element('#duplicate-button')
        remove_button = window.dom.get_element('#remove-button')
        move_button = window.dom.get_element('#move-button')
        empty_button = window.dom.get_element('#empty-button')
        add_button = window.dom.get_element('#add-button')
        color_button = window.dom.get_element('#color-button')
    
        toggle_button.events.click += lambda e: circle.toggle()
        duplicate_button.events.click += lambda e: circle.copy()
        toggle_class_button.events.click += toggle_class
        remove_button.events.click += remove_rectangle
        move_button.events.click += move_rectangle
        empty_button.events.click += empty_container
        add_button.events.click += create_rectangle
        color_button.events.click += change_color
    
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
                            margin: 0.5rem;
                            border-radius: 5px;
                            background-color: red;
                        }
    
                        .circle {
                            border-radius: 50px;
                            background-color: red;
                        }
    
                        .circle:hover {
                            background-color: green;
                        }
    
                        .container {
                            display: flex;
                            flex-wrap: wrap;
                        }
                    </style>
                    </head>
                    <body>
                        <button id="toggle-button">Toggle circle</button>
                        <button id="toggle-class-button">Toggle class</button>
                        <button id="color-button">Change color</button>
                        <button id="duplicate-button">Duplicate circle</button>
                        <button id="add-button">Add rectangle</button>
                        <button id="remove-button" disabled>Remove rectangle</button>
                        <button id="move-button" disabled>Move rectangle</button>
                        <button id="empty-button" disabled>Remove all</button>
                        <div id="circles" style="display: flex; flex-wrap: wrap;">
                            <div id="circle" class="rectangle circle"></div>
                        </div>
    
                        <div id="rectangles" class="container"></div>
                    </body>
                </html>
            '''
        )
        webview.start(bind, window)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/dom_manipulation.md)

[PrevDom Events](/examples/dom_events)[NextDom Traversal](/examples/dom_traversal)
