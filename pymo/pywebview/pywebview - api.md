[Skip to main content](https://pywebview.flowrl.com/api/#main-content)

[![](pywebview - api_files\\logo-no-text.png)pywebview](https://pywebview.flowrl.com/)

[Guide](https://pywebview.flowrl.com/guide/)

[API](https://pywebview.flowrl.com/api/)

[Examples](https://pywebview.flowrl.com/examples/)

[Contributing](https://pywebview.flowrl.com/contributing/)

[Blog](https://pywebview.flowrl.com/blog/)

[Changelog](https://pywebview.flowrl.com/CHANGELOG.html)

[2.x](https://pywebview.flowrl.com/2.4)

[3.x](https://pywebview.flowrl.com/3.7)

[](https://github.com/r0x0r/pywebview)




  1. [pywebview](https://pywebview.flowrl.com/)
  2. [API](https://pywebview.flowrl.com/api/)



# API

October 19, 2018About 11 min

* * *

On This Page

  * [webview.active_window](https://pywebview.flowrl.com/api/#webview-active-window)
  * [webview.create_window](https://pywebview.flowrl.com/api/#webview-create-window)
  * [webview.start](https://pywebview.flowrl.com/api/#webview-start)
  * [webview.screens](https://pywebview.flowrl.com/api/#webview-screens)
  * [webview.settings](https://pywebview.flowrl.com/api/#webview-settings)
  * [webview.token](https://pywebview.flowrl.com/api/#webview-token)
  * [webview.dom](https://pywebview.flowrl.com/api/#webview-dom)
  *     * [webview.dom.DOMEventHandler](https://pywebview.flowrl.com/api/#webview-dom-domeventhandler)
    * [webview.dom.ManipulationMode](https://pywebview.flowrl.com/api/#webview-dom-manipulationmode)
  * [webview.Element](https://pywebview.flowrl.com/api/#webview-element)
  *     * [element.attributes](https://pywebview.flowrl.com/api/#element-attributes)
    * [element.classes](https://pywebview.flowrl.com/api/#element-classes)
    * [element.append](https://pywebview.flowrl.com/api/#element-append)
    * [element.blur](https://pywebview.flowrl.com/api/#element-blur)
    * [element.children](https://pywebview.flowrl.com/api/#element-children)
    * [element.copy](https://pywebview.flowrl.com/api/#element-copy)
    * [element.empty](https://pywebview.flowrl.com/api/#element-empty)
    * [element.events](https://pywebview.flowrl.com/api/#element-events)
    * [element.focus](https://pywebview.flowrl.com/api/#element-focus)
    * [element.focused](https://pywebview.flowrl.com/api/#element-focused)
    * [element.hide](https://pywebview.flowrl.com/api/#element-hide)
    * [element.id](https://pywebview.flowrl.com/api/#element-id)
    * [element.move](https://pywebview.flowrl.com/api/#element-move)
    * [element.next](https://pywebview.flowrl.com/api/#element-next)
    * [element.off](https://pywebview.flowrl.com/api/#element-off)
    * [element.on](https://pywebview.flowrl.com/api/#element-on)
    * [element.parent](https://pywebview.flowrl.com/api/#element-parent)
    * [element.previous](https://pywebview.flowrl.com/api/#element-previous)
    * [element.remove](https://pywebview.flowrl.com/api/#element-remove)
    * [element.show](https://pywebview.flowrl.com/api/#element-show)
    * [element.style](https://pywebview.flowrl.com/api/#element-style)
    * [element.tabindex](https://pywebview.flowrl.com/api/#element-tabindex)
    * [element.tag](https://pywebview.flowrl.com/api/#element-tag)
    * [element.text](https://pywebview.flowrl.com/api/#element-text)
    * [element.toggle](https://pywebview.flowrl.com/api/#element-toggle)
    * [element.value](https://pywebview.flowrl.com/api/#element-value)
    * [element.visible](https://pywebview.flowrl.com/api/#element-visible)
  * [webview.Menu](https://pywebview.flowrl.com/api/#webview-menu)
  *     * [menu.Menu](https://pywebview.flowrl.com/api/#menu-menu)
    * [menu.MenuAction](https://pywebview.flowrl.com/api/#menu-menuaction)
    * [menu.MenuSeparator](https://pywebview.flowrl.com/api/#menu-menuseparator)
  * [webview.Screen](https://pywebview.flowrl.com/api/#webview-screen)
  *     * [screen.height](https://pywebview.flowrl.com/api/#screen-height)
    * [screen.width](https://pywebview.flowrl.com/api/#screen-width)
    * [screen.x](https://pywebview.flowrl.com/api/#screen-x)
    * [screen.y](https://pywebview.flowrl.com/api/#screen-y)
  * [webview.Window](https://pywebview.flowrl.com/api/#webview-window)
  *     * [window.title](https://pywebview.flowrl.com/api/#window-title)
    * [window.on_top](https://pywebview.flowrl.com/api/#window-on-top)
    * [window.x](https://pywebview.flowrl.com/api/#window-x)
    * [window.y](https://pywebview.flowrl.com/api/#window-y)
    * [window.width](https://pywebview.flowrl.com/api/#window-width)
    * [window.height](https://pywebview.flowrl.com/api/#window-height)
    * [window.clear_cookies](https://pywebview.flowrl.com/api/#window-clear-cookies)
    * [window.create_confirmation_dialog](https://pywebview.flowrl.com/api/#window-create-confirmation-dialog)
    * [window.create_file_dialog](https://pywebview.flowrl.com/api/#window-create-file-dialog)
    * [window.destroy](https://pywebview.flowrl.com/api/#window-destroy)
    * [window.evaluate_js](https://pywebview.flowrl.com/api/#window-evaluate-js)
    * [window.expose](https://pywebview.flowrl.com/api/#window-expose)
    * [window.get_cookies](https://pywebview.flowrl.com/api/#window-get-cookies)
    * [window.get_current_url](https://pywebview.flowrl.com/api/#window-get-current-url)
    * [window.get_elements](https://pywebview.flowrl.com/api/#window-get-elements)
    * [window.hide](https://pywebview.flowrl.com/api/#window-hide)
    * [window.load_css](https://pywebview.flowrl.com/api/#window-load-css)
    * [window.load_html](https://pywebview.flowrl.com/api/#window-load-html)
    * [window.load_url](https://pywebview.flowrl.com/api/#window-load-url)
    * [window.maximize](https://pywebview.flowrl.com/api/#window-maximize)
    * [window.minimize](https://pywebview.flowrl.com/api/#window-minimize)
    * [window.move](https://pywebview.flowrl.com/api/#window-move)
    * [window.native](https://pywebview.flowrl.com/api/#window-native)
    * [window.resize](https://pywebview.flowrl.com/api/#window-resize)
    * [window.restore](https://pywebview.flowrl.com/api/#window-restore)
    * [window.run_js](https://pywebview.flowrl.com/api/#window-run-js)
    * [window.set_title](https://pywebview.flowrl.com/api/#window-set-title)
    * [window.show](https://pywebview.flowrl.com/api/#window-show)
    * [window.toggle_fullscreen](https://pywebview.flowrl.com/api/#window-toggle-fullscreen)
    * [window.dom.body](https://pywebview.flowrl.com/api/#window-dom-body)
    * [window.dom.create_element](https://pywebview.flowrl.com/api/#window-dom-create-element)
    * [window.dom.document](https://pywebview.flowrl.com/api/#window-dom-document)
    * [window.dom.get_element](https://pywebview.flowrl.com/api/#window-dom-get-element)
    * [window.dom.get_elements](https://pywebview.flowrl.com/api/#window-dom-get-elements)
    * [window.dom.window](https://pywebview.flowrl.com/api/#window-dom-window)
  * [Window events](https://pywebview.flowrl.com/api/#window-events)
  *     * [window.events.before_show](https://pywebview.flowrl.com/api/#window-events-before-show)
    * [window.events.before_load](https://pywebview.flowrl.com/api/#window-events-before-load)
    * [window.events.closed](https://pywebview.flowrl.com/api/#window-events-closed)
    * [window.events.closing](https://pywebview.flowrl.com/api/#window-events-closing)
    * [window.events.loaded](https://pywebview.flowrl.com/api/#window-events-loaded)
    * [window.events.maximized](https://pywebview.flowrl.com/api/#window-events-maximized)
    * [window.events.minimized](https://pywebview.flowrl.com/api/#window-events-minimized)
    * [window.events.moved](https://pywebview.flowrl.com/api/#window-events-moved)
    * [window.events.restored](https://pywebview.flowrl.com/api/#window-events-restored)
    * [window.events.resized](https://pywebview.flowrl.com/api/#window-events-resized)
    * [window.events.shown](https://pywebview.flowrl.com/api/#window-events-shown)
  * [DOM events](https://pywebview.flowrl.com/api/#dom-events)
  * [Drag area](https://pywebview.flowrl.com/api/#drag-area)



# [API](https://pywebview.flowrl.com/api/#api)

## [webview.active_window](https://pywebview.flowrl.com/api/#webview-active-window)
    
    
    webview.active_window()

Get an instance of the currently active window

## [webview.create_window](https://pywebview.flowrl.com/api/#webview-create-window)
    
    
    webview.create_window(title, url=None, html=None, js_api=None, width=800, height=600,
                          x=None, y=None, screen=None, resizable=True, fullscreen=False,
                          min_size=(200, 100), hidden=False, frameless=False,
                          easy_drag=True, shadow=False, focus=True, minimized=False, maximized=False,
                          on_top=False, confirm_close=False, background_color='#FFFFFF',
                          transparent=False, text_select=False, zoomable=False,
                          draggable=False, server=http.BottleServer, server_args={},
                          localization=None)

Create a new _pywebview_ window and returns its instance. Can be used to create multiple windows (except Android). Window is not shown until the GUI loop is started. If the function is invoked during the GUI loop, the window is displayed immediately.

  * `title` \- Window title
  * `url` \- URL to load. If the URL does not have a protocol prefix, it is resolved as a path relative to the application entry point. Alternatively a WSGI server object can be passed to start a local web server.
  * `html` \- HTML code to load. If both URL and HTML are specified, HTML takes precedence.
  * `js_api` \- Expose a python object to the Javascript domain of the current `pywebview` window. Methods of the `js_api` object can be invoked from Javascript by calling `window.pywebview.api.<methodname>(<parameters>)` functions. Exposed function return a promise that return once function returns. Only basic Python objects (like int, str, dict, ...) can be returned to Javascript.
  * `width` \- Window width. Default is 800px.
  * `height` \- Window height. Default is 600px.
  * `x` \- Window x coordinate. Default is centered.
  * `y` \- Window y coordinate. Default is centered.
  * `screen` \- Screen to display window on. `screen` is a screen instance returned by `webview.screens`.
  * `resizable` \- Whether window can be resized. Default is True
  * `fullscreen` \- Start in fullscreen mode. Default is False
  * `min_size` \- a (width, height) tuple that specifies a minimum window size. Default is 200x100
  * `hidden` \- Create a window hidden by default. Default is False
  * `frameless` \- Create a frameless window. Default is False.
  * `easy_drag` \- Easy drag mode for frameless windows. Window can be moved by dragging any point. Default is True. Note that easy_drag has no effect with normal windows. To control dragging on an element basis, see [drag area](https://pywebview.flowrl.com/guide/api.html#drag-area) for details.
  * `shadow` \- Add window shadow. Default is False. _Windows only_.
  * `focus` \- Create a non-focusable window if False. Default is True.
  * `minimized` \- Display window minimized
  * `maximized` \- Display window maximized
  * `on_top` \- Set window to be always on top of other windows. Default is False.
  * `confirm_close` \- Whether to display a window close confirmation dialog. Default is False
  * `background_color` \- Background color of the window displayed before WebView is loaded. Specified as a hex color. Default is white.
  * `transparent` \- Create a transparent window. Not supported on Windows. Default is False. Note that this setting does not hide or make window chrome transparent. To hide window chrome set `frameless` to True.
  * `text_select` \- Enables document text selection. Default is False. To control text selection on per element basis, use [user-select](https://developer.mozilla.org/en-US/docs/Web/CSS/user-select) CSS property.
  * `zoomable` \- Enable document zooming. Default is False
  * `draggable` \- Enable image and link object dragging. Default is False server=http.BottleServer, server_args
  * `vibrancy` \- Enable window vibrancy. Default is False. macOS only.
  * `server` \- A custom WSGI server instance for this window. Defaults to BottleServer.
  * `server_args` \- Dictionary of arguments to pass through to the server instantiation
  * `localization` \- pass a localization dictionary for per window localization.



## [webview.start](https://pywebview.flowrl.com/api/#webview-start)
    
    
    webview.start(func=None, args=None, localization={}, gui=None, debug=False,
                  http_server=False, http_port=None, user_agent=None, private_mode=True,
                  storage_path=None, menu=[], server=http.BottleServer, ssl=False,
                  server_args={}, icon=None):

Start a GUI loop and display previously created windows. This function must be called from a main thread.

  * `func` \- function to invoke upon starting the GUI loop.
  * `args` \- function arguments. Can be either a single value or a tuple of values.
  * `localization` \- a dictionary with localized strings. Default strings and their keys are defined in localization.py
  * `gui` \- force a specific GUI. Allowed values are `cef`, `qt` or `gtk` depending on a platform. See [Web Engine](https://pywebview.flowrl.com/guide/web_engine.html) for details.
  * `debug` \- enable debug mode. See [Debugging](https://pywebview.flowrl.com/guide/debugging.html) for details.
  * `http_server` \- enable built-in HTTP server for absolute local paths. For relative paths HTTP server is started automatically and cannot be disabled. For each window, a separate HTTP server is spawned. This option is ignored for non-local URLs.
  * `http_port` \- specify a port number for the HTTP server. By default port is randomized.
  * `user_agent` \- change user agent string.
  * `private_mode` \- Control whether cookies and other persistant objects are stored between session. By default private mode is on and nothing is stored between sessions.
  * `storage_path` \- An optional location on hard drive where to store persistant objects like cookies and local storage. By default `~/.pywebview` is used on *nix systems and `%APPDATA%\pywebview` on Windows.
  * `menu` \- Pass a list of Menu objects to create an application menu. See [this example](https://pywebview.flowrl.com/examples/menu.html) for usage details.
  * `server` \- A custom WSGI server instance. Defaults to BottleServer.
  * `ssl` \- If using the default BottleServer (and for now the GTK backend), will use SSL encryption between the webview and the internal server. You need to have `cryptography` pip dependency installed in order to use `ssl`. It is not installed by default.
  * `server_args` \- Dictionary of arguments to pass through to the server instantiation
  * `icon` \- path to application icon. Available only for GTK / QT. For other platforms icon should be specified via a bundler.



#### [Examples](https://pywebview.flowrl.com/api/#examples)

  * [Simple window](https://pywebview.flowrl.com/examples/open_url.html)
  * [Multi-window](https://pywebview.flowrl.com/examples/multiple_windows.html)



## [webview.screens](https://pywebview.flowrl.com/api/#webview-screens)
    
    
    webview.screens

Return a list of available displays (as `Screen` objects) with the primary display as the first element of the list.

#### [Examples](https://pywebview.flowrl.com/api/#examples-1)

  * [Simple window](https://pywebview.flowrl.com/examples/screens.html)



## [webview.settings](https://pywebview.flowrl.com/api/#webview-settings)
    
    
    webview.settings = {
      'ALLOW_DOWNLOADS': False,
      'ALLOW_FILE_URLS': True,
      'OPEN_EXTERNAL_LINKS_IN_BROWSER': True,
      'OPEN_DEVTOOLS_IN_DEBUG': True,
      'REMOTE_DEBUGGING_PORT': None
    }

Additional options that override default behaviour of _pywebview_ to address popular feature requests.

  * `ALLOW_DOWNLOADS` Allow file downloads. Disabled by default.
  * `ALLOW_FILE_URLS` Enable `file://` urls. Disabled by default.
  * `OPEN_EXTERNAL_LINKS_IN_BROWSER`. Open `target=_blank` link in an external browser. Enabled by default.
  * `OPEN_DEVTOOLS_IN_DEBUG` Open devtools automatically in debug mode. Enabled by default.
  * `REMOTE_DEBUGGING_PORT` Enable remote debugging when using `edgechromium`. Disabled by default.



#### [Examples](https://pywebview.flowrl.com/api/#examples-2)

  * [File downloads](https://pywebview.flowrl.com/examples/downloads.html)



## [webview.token](https://pywebview.flowrl.com/api/#webview-token)
    
    
    webview.token

A CSRF token property unique to the session. The same token is exposed as `window.pywebview.token`. See [Security](https://pywebview.flowrl.com/guide/security.html) for usage details.

## [webview.dom](https://pywebview.flowrl.com/api/#webview-dom)

### [webview.dom.DOMEventHandler](https://pywebview.flowrl.com/api/#webview-dom-domeventhandler)
    
    
    DOMEventHandler(callback, prevent_default=False, stop_propagation=False, stop_immediate_propagation=False, debounce=0)

A container for an event handler used to control propagation or default behaviour of the event. If `debounce` is greater than zero, Python event handler is debounced by a specified number of milliseconds. This can be useful for events like `dragover` and `mouseover` that generate a constant stream of events resulting in poor performance.

#### [Examples](https://pywebview.flowrl.com/api/#examples-3)
    
    
    element.events.click += DOMEventHandler(on_click, prevent_default=True, stop_propagation=True, stop_immediate_propagation=True)
    element.events.mouseover += DOMEventHandler(on_click, debounce=500)

### [webview.dom.ManipulationMode](https://pywebview.flowrl.com/api/#webview-dom-manipulationmode)

Enum that sets the position of a manipulated DOM element. Possible values are:

  * `LastChild` \- element is inserted as a last child of the target
  * `FirstChild` \- element is inserted as a firt child of the target
  * `Before` \- element is inserted before the target
  * `After` \- element is inserted after the target
  * `Replace` \- element is inserted replacing the target



Used by `element.append`, `element.copy`, `element.move` and `window.dom.create_element` functions.

## [webview.Element](https://pywebview.flowrl.com/api/#webview-element)

### [element.attributes](https://pywebview.flowrl.com/api/#element-attributes)

Get or modify element's attributes. `attributes` is a `PropsDict` dict-like object that implements most of dict functions. To add an attribute, you can simply assign a value to a key in `attributes`. Similarly, to remove an attribute, you can set its value to None.

#### [Examples](https://pywebview.flowrl.com/api/#examples-4)
    
    
    element.attributes['id'] = 'container-id' # set element's id
    element.attributes['data-flag'] = '1337'
    element.attributes['id'] = None # remove element's id
    del element.attributes['data-flag'] # remove element's data-flag attribute

### [element.classes](https://pywebview.flowrl.com/api/#element-classes)
    
    
    element.classes

Get or set element's classes. `classes` is a `ClassList` list-like object that implements a subset of list functions like `append`, `remove` and `clear`. Additionally it has a `toggle` function for toggling a class.

#### [Examples](https://pywebview.flowrl.com/api/#examples-5)
    
    
    element.classes = ['container', 'red', 'dotted'] # overwrite element's classes
    element.classes.remove('red') # remove red class
    element.classes.add('blue') # add blue class
    element.classes.toggle('dotted')

### [element.append](https://pywebview.flowrl.com/api/#element-append)
    
    
    element.append(html, mode=webview.dom.ManipulationMode.LastChild)

Insert HTML content to the element as a last child. To control the position of the new element, use the `mode` parameter. See [Manipulation mode](https://pywebview.flowrl.com/guide/api.html#manipulation-mode) for possible values.

### [element.blur](https://pywebview.flowrl.com/api/#element-blur)
    
    
    element.blur()

Blur element.

### [element.children](https://pywebview.flowrl.com/api/#element-children)
    
    
    element.children

Get element's children elements. Returns a list of `Element` objects.

### [element.copy](https://pywebview.flowrl.com/api/#element-copy)
    
    
    element.copy(target=None, mode=webview.dom.ManipulationMode.LastChild, id=None)

Create a new copy of the element. `target` can be either another `Element` or a DOM selector string. If target is omitted, a copy is created in the current element's parent. To control the position of the new element, use the `mode` parameter. See [Manipulation mode](https://pywebview.flowrl.com/guide/api.html#manipulation-mode) for possible values. The id parameter is stripped from the copy. Optionally you can set the id of the copy by specifying the `id` parameter.

### [element.empty](https://pywebview.flowrl.com/api/#element-empty)
    
    
    element.empty()

Empty element by removing all its children.

### [element.events](https://pywebview.flowrl.com/api/#element-events)
    
    
    element.events

A container object of element's all DOM events, ie `events.click`, `event.keydown`. This container is dynamically populated and its contents depend on the events a node has. To subscribe to a DOM event, use the `+=` syntax, e.g. `element.events.click += callback`. Similarly to remove an event listener use `-=`, eg. `element.events.click -= callback`. Callback can be either a function or an instance of `DOMEventHandler` if you need to control propagation of the event.

### [element.focus](https://pywebview.flowrl.com/api/#element-focus)
    
    
    element.focus()

Focus element.

### [element.focused](https://pywebview.flowrl.com/api/#element-focused)
    
    
    element.focused

Get whether the element is focused.

### [element.hide](https://pywebview.flowrl.com/api/#element-hide)
    
    
    element.hide()

Hide element by setting `display: none`.

### [element.id](https://pywebview.flowrl.com/api/#element-id)
    
    
    element.id

Get or set element's id. None if id is not set.

### [element.move](https://pywebview.flowrl.com/api/#element-move)
    
    
    element.move(target, mode=webview.dom.ManipulationMode.LastChild)

Move element to the `target` that can be either another `Element` or a DOM selector string. To control the position of the new element, use the `mode` parameter. See [Manipulation mode](https://pywebview.flowrl.com/guide/api.html#manipulation-mode) for possible values.

#### [Examples](https://pywebview.flowrl.com/api/#examples-6)

[DOM Manipulation](https://pywebview.flowrl.com/examples/dom_manipulation.html)

### [element.next](https://pywebview.flowrl.com/api/#element-next)
    
    
    element.next

Get element's next sibling. None if no sibling is present.

### [element.off](https://pywebview.flowrl.com/api/#element-off)
    
    
    element.off(event, callback)

Remove an event listener. Identical to `element.event.event_name -= callback`.

#### [Examples](https://pywebview.flowrl.com/api/#examples-7)
    
    
    # these two are identical
    element.off('click', callback_func)
    element.events.click -= callback_func

[DOM Events](https://pywebview.flowrl.com/examples/dom_events.html)

### [element.on](https://pywebview.flowrl.com/api/#element-on)
    
    
    element.on(event, callback)

Add an event listener to a DOM event. Callback can be either a function or an instance of `DOMEventHandler` if you need to control propagation of the event. Identical to `element.event.event_name += callback`.

#### [Examples](https://pywebview.flowrl.com/api/#examples-8)
    
    
    # these two are identical
    element.on('click', callback_func)
    element.events.click += callback_func

[DOM Events](https://pywebview.flowrl.com/examples/dom_events.html)

### [element.parent](https://pywebview.flowrl.com/api/#element-parent)
    
    
    element.parent

Get element's parent `Element` or None if root element is reached.

### [element.previous](https://pywebview.flowrl.com/api/#element-previous)
    
    
    element.previous

Get element's previous sibling. None if no sibling is present.

### [element.remove](https://pywebview.flowrl.com/api/#element-remove)
    
    
    element.remove()

Remove element from DOM. `Element` object is not destroyed, but marked as removed. Trying to access any properties or invoke any functions of the element will result in a warning.

[DOM Manipulation](https://pywebview.flowrl.com/examples/dom_manipulation.html)

### [element.show](https://pywebview.flowrl.com/api/#element-show)
    
    
    element.show()

Show hidden element. If element was hidden with `element.hide()`, a previous display value is restored. Otherwise `display: block` is set.

[DOM Manipulation](https://pywebview.flowrl.com/examples/dom_manipulation.html)

### [element.style](https://pywebview.flowrl.com/api/#element-style)

Get or modify element's styles. `style` is a `PropsDict` dict-like object that implements most of dict functions. To add a style declraration, you can simply assign a value to a key in `attributes`. Similarly, to reset a declaration, you can set its value to None.

#### [Examples](https://pywebview.flowrl.com/api/#examples-9)
    
    
    element.style['width'] = '100px' # set element's width to 100px
    element.style['display'] = 'flex' # set element's display property to flex
    element.style['width'] = None # reset width to auto
    del element.attributes['display'] # reset display property to block

### [element.tabindex](https://pywebview.flowrl.com/api/#element-tabindex)
    
    
    element.tabindex

Get or set element's tabindex.

### [element.tag](https://pywebview.flowrl.com/api/#element-tag)
    
    
    element.tag

Get element's tag name.

### [element.text](https://pywebview.flowrl.com/api/#element-text)
    
    
    element.text

Get or set element's text content.

### [element.toggle](https://pywebview.flowrl.com/api/#element-toggle)
    
    
    element.toggle()

Toggle element's visibility.

### [element.value](https://pywebview.flowrl.com/api/#element-value)
    
    
    element.value

Get or set element's value. Applicable only to input elements that have a value.

### [element.visible](https://pywebview.flowrl.com/api/#element-visible)
    
    
    element.visible

Get whether the element is visible.

## [webview.Menu](https://pywebview.flowrl.com/api/#webview-menu)

Used to create an application menu. See [this example](https://pywebview.flowrl.com/examples/menu.html) for usage details.

### [menu.Menu](https://pywebview.flowrl.com/api/#menu-menu)

`Menu(title, items=[])`. Instantiate to create a menu that can be either top level menu or a nested menu. `title` is the title of the menu and `items` is a list of actions, separators or other menus.

### [menu.MenuAction](https://pywebview.flowrl.com/api/#menu-menuaction)

`MenuAction(title, function)` Instantiate to create a menu item. `title` is the name of the item and function is a callback that should be called when menu action is clicked.

### [menu.MenuSeparator](https://pywebview.flowrl.com/api/#menu-menuseparator)

`MenuSeparator()` Instantiate to create a menu separator.

## [webview.Screen](https://pywebview.flowrl.com/api/#webview-screen)

Represents a display found on the systems. A list of `Screen` objects is returned by `webview.screens` property.

### [screen.height](https://pywebview.flowrl.com/api/#screen-height)
    
    
    screen.height

Get display height.

### [screen.width](https://pywebview.flowrl.com/api/#screen-width)
    
    
    screen.width

### [screen.x](https://pywebview.flowrl.com/api/#screen-x)
    
    
    screen.x

### [screen.y](https://pywebview.flowrl.com/api/#screen-y)
    
    
    screen.y

Get display width.

## [webview.Window](https://pywebview.flowrl.com/api/#webview-window)

Represents a window that hosts webview. `window` object is returned by `create_window` function.

### [window.title](https://pywebview.flowrl.com/api/#window-title)
    
    
    window.title

Get or set title of the window.

### [window.on_top](https://pywebview.flowrl.com/api/#window-on-top)
    
    
    window.on_top

Get or set whether the window is always on top.

### [window.x](https://pywebview.flowrl.com/api/#window-x)
    
    
    window.x

Get X coordinate of the top-left corrner of the window.

### [window.y](https://pywebview.flowrl.com/api/#window-y)
    
    
    window.y

Get Y coordinate of the top-left corrner of the window.

### [window.width](https://pywebview.flowrl.com/api/#window-width)
    
    
    window.width

Get width of the window

### [window.height](https://pywebview.flowrl.com/api/#window-height)
    
    
    window.height

Get height of the window

### [window.clear_cookies](https://pywebview.flowrl.com/api/#window-clear-cookies)
    
    
    window.clear_cookies()

Clear all the cookies including `HttpOnly` ones.

#### [Example](https://pywebview.flowrl.com/api/#example)

  * [Cookies](https://pywebview.flowrl.com/examples/cookies.html)



### [window.create_confirmation_dialog](https://pywebview.flowrl.com/api/#window-create-confirmation-dialog)
    
    
    window.create_confirmation_dialog(title, message)

Create a confirmation (Ok / Cancel) dialog.

### [window.create_file_dialog](https://pywebview.flowrl.com/api/#window-create-file-dialog)
    
    
    window.create_file_dialog(dialog_type=OPEN_DIALOG, directory='', allow_multiple=False, save_filename='', file_types=())

Create an open file (`webview.OPEN_DIALOG`), open folder (`webview.FOLDER_DIALOG`) or save file (`webview.SAVE_DIALOG`) dialog.

Return a tuple of selected files, None if cancelled.

  * `allow_multiple=True` enables multiple selection.
  * `directory` Initial directory.
  * `save_filename` Default filename for save file dialog.
  * `file_types` A tuple of supported file type strings in the open file dialog. A file type string must follow this format `"Description (*.ext1;*.ext2...)"`.



If the argument is not specified, then the `"All files (*.*)"` mask is used by default. The 'All files' string can be changed in the localization dictionary.

#### [Examples](https://pywebview.flowrl.com/api/#examples-10)

  * [Open-file dialog](https://pywebview.flowrl.com/examples/open_file_dialog.html)
  * [Save-file dialog](https://pywebview.flowrl.com/examples/save_file_dialog.html)



### [window.destroy](https://pywebview.flowrl.com/api/#window-destroy)
    
    
    window.destroy()

Destroy the window.

[Example](https://pywebview.flowrl.com/examples/destroy_window.html)

### [window.evaluate_js](https://pywebview.flowrl.com/api/#window-evaluate-js)
    
    
    window.evaluate_js(script, callback=None)

Execute Javascript code. The last evaluated expression is returned. If callback function is supplied, then promises are resolved and the callback function is called with the result as a parameter. Javascript types are converted to Python types, eg. JS objects to dicts, arrays to lists, undefined to None. DOM nodes are serialized using custom serialization. Functions are omitted and circular references are converted to the `[Circular Reference]` string literal. `webview.error.JavascriptException` is thrown if executed codes raises an error. r-strings is a recommended way to load Javascript. Note that the `evaluate_js` employs `eval`, which will fail if `unsafe-eval` CSP is set. Alternatively you may use `window.run_js(code)` that executes Javascript code as is without returning a result.

### [window.expose](https://pywebview.flowrl.com/api/#window-expose)

Expose a Python function or functions to JS API. Functions are exposed as `window.pywebview.api.func_name`

[Example](https://pywebview.flowrl.com/examples/expose.html)

### [window.get_cookies](https://pywebview.flowrl.com/api/#window-get-cookies)
    
    
    window.get_cookies()

Return a list of all the cookies set for the current website (as [SimpleCookie](https://docs.python.org/3/library/http.cookies.html)).

### [window.get_current_url](https://pywebview.flowrl.com/api/#window-get-current-url)
    
    
    window.get_current_url()

Return the current URL. None if no url is loaded.

[Example](https://pywebview.flowrl.com/examples/get_current_url.html)

### [window.get_elements](https://pywebview.flowrl.com/api/#window-get-elements)
    
    
    window.get_elements(selector)

_DEPRECATED_. Use `window.dom.get_elements` instead.

[Example](https://pywebview.flowrl.com/examples/get_elements.html)

### [window.hide](https://pywebview.flowrl.com/api/#window-hide)
    
    
    window.hide()

Hide the window.

[Example](https://pywebview.flowrl.com/examples/show_hide.html)

### [window.load_css](https://pywebview.flowrl.com/api/#window-load-css)
    
    
    window.load_css(css)

Load CSS as a string.

[Example](https://pywebview.flowrl.com/examples/css_load.html)

### [window.load_html](https://pywebview.flowrl.com/api/#window-load-html)
    
    
    window.load_html(content, base_uri=base_uri())

Load HTML code. Base URL for resolving relative URLs is set to the directory the program is launched from. Note that you cannot use hashbang anchors when HTML is loaded this way.

[Example](https://pywebview.flowrl.com/examples/html_load.html)

### [window.load_url](https://pywebview.flowrl.com/api/#window-load-url)
    
    
    window.load_url(url)

Load a new URL.

[Example](https://pywebview.flowrl.com/examples/change_url.html)

### [window.maximize](https://pywebview.flowrl.com/api/#window-maximize)
    
    
    window.maximize()

Maximize window.

[Example](https://pywebview.flowrl.com/examples/window_state.html)

### [window.minimize](https://pywebview.flowrl.com/api/#window-minimize)
    
    
    window.minimize()

Minimize window.

[Example](https://pywebview.flowrl.com/examples/window_state.html)

### [window.move](https://pywebview.flowrl.com/api/#window-move)
    
    
    window.move(x, y)

Move window to a new position.

[Example](https://pywebview.flowrl.com/examples/move_window.html)

### [window.native](https://pywebview.flowrl.com/api/#window-native)
    
    
    window.native.Handle # get application window handle on Windows

Get a native window object. This can be useful for applying custom styling to the window. Object type depends on the platform

`System.Windows.Form` \- Windows `AppKit.NSWindow` \- macOS `Gtk.ApplicationWindow` \- GTK `QMainWindow` \- QT `kivy.uix.widget.Widget` \- Android

The `native` property is available after the `before_show` event is fired.

You can also each platform's WebView object via `window.native.webview`. WebView's types are as follows.

`Microsoft.Web.WebView2.WinForms.WebView2` \- Windows / EdgeChromium `System.Windows.Forms.WebBrowser` \- Windows / MSHTML `WebKit.WKWebView` \- macOS `gi.repository.WebKit2.WebView` \- GTK `QtWebEngineWidgets.QWebEngineView` / `QtWebKitWidgets.QWebView`\- QT `android.webkit.WebView` \- Android

### [window.resize](https://pywebview.flowrl.com/api/#window-resize)
    
    
    window.resize(width, height, fix_point=FixPoint.NORTH | FixPoint.WEST)

Resize window. Optional parameter fix_point specifies in respect to which point the window is resized. The parameter accepts values of the `webview.window.FixPoint` enum (`NORTH`, `SOUTH`, `EAST`, `WEST`)

[Example](https://pywebview.flowrl.com/examples/minimize.html)

### [window.restore](https://pywebview.flowrl.com/api/#window-restore)
    
    
    window.restore()

Restore minimized window.

[Example](https://pywebview.flowrl.com/examples/minimize.html)

### [window.run_js](https://pywebview.flowrl.com/api/#window-run-js)
    
    
    window.run_js('document.body.style.color = "deepred"')

Execute Javascript as is without wrapping it in `eval` and helper code. This function does not return a result.

[Example](https://pywebview.flowrl.com/examples/run_js.html)

### [window.set_title](https://pywebview.flowrl.com/api/#window-set-title)
    
    
    window.set_title(title)

_DEPRECATED_. Use `window.title` instead. Change the title of the window.

[Example](https://pywebview.flowrl.com/examples/window_title_change.html)

### [window.show](https://pywebview.flowrl.com/api/#window-show)
    
    
    window.show()

Show the window if it is hidden. Has no effect otherwise

[Example](https://pywebview.flowrl.com/examples/show_hide.html)

### [window.toggle_fullscreen](https://pywebview.flowrl.com/api/#window-toggle-fullscreen)
    
    
    window.toggle_fullscreen()

Toggle fullscreen mode on the active monitor.

[Example](https://pywebview.flowrl.com/examples/toggle_fullscreen.html)

### [window.dom.body](https://pywebview.flowrl.com/api/#window-dom-body)
    
    
    window.body

Get document's body as an `Element` object

### [window.dom.create_element](https://pywebview.flowrl.com/api/#window-dom-create-element)
    
    
    window.create_element(html, parent=None, mode=webview.dom.ManipulationMode.LastChild)

Insert HTML content and returns the Element of the root object. `parent` can be either another `Element` or a DOM selector string. If parent is omited, created DOM is attached to document's body. To control the position of the new element, use the `mode` parameter. See [Manipulation mode](https://pywebview.flowrl.com/guide/api.html#manipulation-mode) for possible values.

### [window.dom.document](https://pywebview.flowrl.com/api/#window-dom-document)
    
    
    window.document

Get `window.document` of the loaded page as an `Element` object

### [window.dom.get_element](https://pywebview.flowrl.com/api/#window-dom-get-element)
    
    
    window.get_element(selector: str)

Get a first `Element` matching the selector. None if not found.

### [window.dom.get_elements](https://pywebview.flowrl.com/api/#window-dom-get-elements)
    
    
    window.get_elements(selector: str)

Get a list of `Element` objects matching the selector.

### [window.dom.window](https://pywebview.flowrl.com/api/#window-dom-window)

Get DOM document's window `window` as an `Element` object

## [Window events](https://pywebview.flowrl.com/api/#window-events)

Window object exposes various lifecycle and window management events. To subscribe to an event, use the `+=` syntax, e.g., `window.events.loaded += func`. Duplicate subscriptions are ignored, and the function is invoked only once for duplicate subscribers. To unsubscribe, use the `-=` syntax, e.g., `window.events.loaded -= func`. To access the window object from the event handler, supply the `window` parameter as the first positional argument of the handler. Most window events are asynchronous, and event handlers are executed in separate threads. The `before_show` and `before_load` events are synchronous and block the main thread until handled.

### [window.events.before_show](https://pywebview.flowrl.com/api/#window-events-before-show)

This event is fired just before pywebview window is shown. This is the earliest event that exposes `window.native` property. This event is blocking.

### [window.events.before_load](https://pywebview.flowrl.com/api/#window-events-before-load)

The event is fired right before _pywebview_ code is injected into the page. The event roughly corresponds to `DOMContentLoaded` DOM event. This event is blocking.

### [window.events.closed](https://pywebview.flowrl.com/api/#window-events-closed)

The event is fired just before _pywebview_ window is closed.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.closing](https://pywebview.flowrl.com/api/#window-events-closing)

The event is fired when _pywebview_ window is about to be closed. If confirm_close is set, then this event is fired before the close confirmation is displayed. If event handler returns False, the close operation will be cancelled.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.loaded](https://pywebview.flowrl.com/api/#window-events-loaded)

The event is fired when DOM is ready.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.maximized](https://pywebview.flowrl.com/api/#window-events-maximized)

The event is fired when window is maximized (fullscreen on macOS)

### [window.events.minimized](https://pywebview.flowrl.com/api/#window-events-minimized)

The event is fired when window is minimized.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.moved](https://pywebview.flowrl.com/api/#window-events-moved)

The event is fired when window is moved.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.restored](https://pywebview.flowrl.com/api/#window-events-restored)

The event is fired when window is restored.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.resized](https://pywebview.flowrl.com/api/#window-events-resized)

The event is fired when pywebview window is resized. Event handler can either have no or accept (width, height) arguments.

[Example](https://pywebview.flowrl.com/examples/events.html)

### [window.events.shown](https://pywebview.flowrl.com/api/#window-events-shown)

The event is fired when pywebview window is shown.

[Example](https://pywebview.flowrl.com/examples/events.html)

## [DOM events](https://pywebview.flowrl.com/api/#dom-events)

 _pywebview_ exposes a `window.pywebviewready` DOM event that is fired after `window.pywebview` is created.

[Example](https://pywebview.flowrl.com/examples/js_api.html)

## [Drag area](https://pywebview.flowrl.com/api/#drag-area)

With a frameless _pywebview_ window, A window can be moved or dragged by adding a special class called `pywebview-drag-region` to any element.
    
    
    <div class='pywebview-drag-region'>Now window can be moved by dragging this DIV.</div>

The magic class name can be overriden by re-assigning the `webview.DRAG_REGION_SELECTOR` constant.

[Example](https://pywebview.flowrl.com/examples/drag_region.html)

[Edit this page](https://github.com/r0x0r/pywebview/edit/docs/docs/api/README.md)

Last update: 2/11/2025, 3:13:14 PM

↻
