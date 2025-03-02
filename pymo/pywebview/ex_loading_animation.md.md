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



# Loading Animation

May 17, 2017About 1 min

* * *

# Loading Animation

Create a loading animation that is displayed before application is loaded.
    
    
    import threading
    
    import webview
    
    html = """
        <style>
            body {
                background-color: #333;
                color: white;
                font-family: Helvetica Neue, Helvetica, Arial, sans-serif;
            }
    
            .main-container {
                width: 100%;
                height: 90vh;
                display: flex;
                display: -webkit-flex;
                align-items: center;
                -webkit-align-items: center;
                justify-content: center;
                -webkit-justify-content: center;
                overflow: hidden;
            }
    
            .loading-container {
            }
    
            .loader {
              font-size: 10px;
              margin: 50px auto;
              text-indent: -9999em;
              width: 3rem;
              height: 3rem;
              border-radius: 50%;
              background: #ffffff;
              background: -moz-linear-gradient(left, #ffffff 10%, rgba(255, 255, 255, 0) 42%);
              background: -webkit-linear-gradient(left, #ffffff 10%, rgba(255, 255, 255, 0) 42%);
              background: -o-linear-gradient(left, #ffffff 10%, rgba(255, 255, 255, 0) 42%);
              background: -ms-linear-gradient(left, #ffffff 10%, rgba(255, 255, 255, 0) 42%);
              background: linear-gradient(to right, #ffffff 10%, rgba(255, 255, 255, 0) 42%);
              position: relative;
              -webkit-animation: load3 1.4s infinite linear;
              animation: load3 1.4s infinite linear;
              -webkit-transform: translateZ(0);
              -ms-transform: translateZ(0);
              transform: translateZ(0);
            }
            .loader:before {
              width: 50%;
              height: 50%;
              background: #ffffff;
              border-radius: 100% 0 0 0;
              position: absolute;
              top: 0;
              left: 0;
              content: '';
            }
            .loader:after {
              background: #333;
              width: 75%;
              height: 75%;
              border-radius: 50%;
              content: '';
              margin: auto;
              position: absolute;
              top: 0;
              left: 0;
              bottom: 0;
              right: 0;
            }
            @-webkit-keyframes load3 {
              0% {
                -webkit-transform: rotate(0deg);
                transform: rotate(0deg);
              }
              100% {
                -webkit-transform: rotate(360deg);
                transform: rotate(360deg);
              }
            }
            @keyframes load3 {
              0% {
                -webkit-transform: rotate(0deg);
                transform: rotate(0deg);
              }
              100% {
                -webkit-transform: rotate(360deg);
                transform: rotate(360deg);
              }
            }
    
            .loaded-container {
                display: none;
            }
    
        </style>
        <body>
          <div class="main-container">
              <div id="loader" class="loading-container">
                  <div class="loader">Loading...</div>
              </div>
    
              <div id="main" class="loaded-container">
                  <h1>Content is loaded!</h1>
              </div>
          </div>
    
          <script>
              setTimeout(function() {
                  document.getElementById('loader').style.display = 'none'
                  document.getElementById('main').style.display = 'block'
              }, 5000)
          </script>
        </body>
    """
    
    if __name__ == '__main__':
        window = webview.create_window('Loading Animation', html=html, background_color='#333333')
        webview.start()

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/examples/loading_animation.md)

Last update: 9/26/2023, 6:21:55 AM

[PrevLoad Html](/examples/load_html)[NextLocalhost Ssl](/examples/localhost_ssl)
