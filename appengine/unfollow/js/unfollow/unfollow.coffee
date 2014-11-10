debug.setLevel(9) if window.location.search is '?debug'

((window) ->
  'use strict'

  namespace = (target, name, block) ->
    [target, name, block] = [window, arguments...] if arguments.length < 3
    top = target
    target = target[item] or= {} for item in name.split '.'
    block target, top

  namespace 'unfollow', (exports, top) ->
    exports.namespace = namespace

    callbacks = []
    exports.init = (callback) ->
      callbacks.push callback

    exports.run = ->
      callback() for callback in callbacks

    $ ->
      unfollow.run()

)(window)

unfollow.namespace 'unfollow.settings', (exports) ->
  'use strict'

  load = ->
    $('div.js-settings div').each ->
      key = $(@).data('key')
      value =  $(@).data('value')
      return if not key
      exports[key] = value
      debug.setLevel(9) if key is 'DEBUG' and value
      debug.info 'unfollow.settings.' + key, value
    $('div.js-settings').remove()

  config = ->
    require.config
      baseUrl: unfollow.settings.STATIC + '/' + unfollow.settings.VERSION + '/lib'
      paths:
        ga: ['//www.google-analytics.com/analytics']

  load()
  config()

$ ->
  # these are executed only once on page load

  $.pjax.defaults.success = -> unfollow.run()
  $('a').pjax
    containers: ['#base']

