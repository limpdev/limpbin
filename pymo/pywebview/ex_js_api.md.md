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



# Js Api

November 7, 2017About 1 min

* * *

# Js Api

Create an application without a HTTP server. The application uses Javascript API object to communicate between Python and Javascript.
    
    
    import random
    import sys
    import threading
    import time
    
    import webview
    
    html = """
    <!DOCTYPE html>
    <html>
    <head lang="en">
    <meta charset="UTF-8">
    
    <style>
        #response-container {
            display: none;
            padding: 1rem;
            margin: 3rem 5%;
            font-size: 120%;
            border: 5px dashed #ccc;
            word-wrap: break-word;
        }
    
        label {
            margin-left: 0.3rem;
            margin-right: 0.3rem;
        }
    
        button {
            font-size: 100%;
            padding: 0.5rem;
            margin: 0.3rem;
            text-transform: uppercase;
        }
    
    </style>
    </head>
    <body>
    
    <h1>JS API Example</h1>
    <p id='pywebview-status'><i>pywebview</i> is not ready</p>
    
    <button onClick="initialize()">Hello Python</button><br/>
    <button id="heavy-stuff-btn" onClick="doHeavyStuff()">Perform a heavy operation</button><br/>
    <button onClick="getRandomNumber()">Get a random number</button><br/>
    <label for="name_input">Say hello to:</label><input id="name_input" placeholder="put a name here">
    <button onClick="greet()">Greet</button><br/>
    <button onClick="catchException()">Catch Exception</button><br/>
    
    <div id="response-container"></div>
    <script>
        window.addEventListener('pywebviewready', function() {
            var container = document.getElementById('pywebview-status')
            container.innerHTML = '<i>pywebview</i> is ready'
        })
    
        function showResponse(response) {
            var container = document.getElementById('response-container')
    
            container.innerText = response.message
            container.style.display = 'block'
        }
    
        function initialize() {
            pywebview.api.init().then(showResponse)
        }
    
        function doHeavyStuff() {
            var btn = document.getElementById('heavy-stuff-btn')
    
            pywebview.api.heavy_stuff.doHeavyStuff().then(function(response) {
                showResponse(response)
                btn.onclick = doHeavyStuff
                btn.innerText = 'Perform a heavy operation'
            })
    
            showResponse({message: 'Working...'})
            btn.innerText = 'Cancel the heavy operation'
            btn.onclick = cancelHeavyStuff
        }
    
        function cancelHeavyStuff() {
            pywebview.api.heavy_stuff.cancelHeavyStuff()
        }
    
        function getRandomNumber() {
            pywebview.api.getRandomNumber().then(showResponse)
        }
    
        function greet() {
            var name_input = document.getElementById('name_input').value;
            pywebview.api.sayHelloTo(name_input).then(showResponse)
        }
    
        function catchException() {
            pywebview.api.error().catch(showResponse)
        }
    
    </script>
    </body>
    </html>
    """
    
    class HeavyStuffAPI:
        def __init__(self):
            self.cancel_heavy_stuff_flag = False
    
        def doHeavyStuff(self):
            time.sleep(0.1)  # sleep to prevent from the ui thread from freezing for a moment
            now = time.time()
            self.cancel_heavy_stuff_flag = False
            for i in range(0, 1000000):
                _ = i * random.randint(0, 1000)
                if self.cancel_heavy_stuff_flag:
                    response = {'message': 'Operation cancelled'}
                    break
            else:
                then = time.time()
                response = {
                    'message': 'Operation took {0:.1f} seconds on the thread {1}'.format(
                        (then - now), threading.current_thread()
                    )
                }
            return response
    
        def cancelHeavyStuff(self):
            time.sleep(0.1)
            self.cancel_heavy_stuff_flag = True
    
    class NotExposedApi:
        def notExposedMethod(self):
            return 'This method is not exposed'
    
    class Api:
        heavy_stuff = HeavyStuffAPI()
        _this_wont_be_exposed = NotExposedApi()
    
        def init(self):
            response = {'message': 'Hello from Python {0}'.format(sys.version)}
            return response
    
        def getRandomNumber(self):
            response = {
                'message': 'Here is a random number courtesy of randint: {0}'.format(
                    random.randint(0, 100000000)
                )
            }
            return response
    
        def sayHelloTo(self, name):
            response = {'message': 'Hello {0}!'.format(name)}
            return response
    
        def error(self):
            raise Exception('This is a Python exception')
    
    if __name__ == '__main__':
        api = Api()
        window = webview.create_window('JS API example', html=html, js_api=api)
        webview.start()

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/js_api.md)

Last update: 9/26/2023, 6:21:55 AM

[PrevIcon](/examples/icon)[NextLinks](/examples/links)
