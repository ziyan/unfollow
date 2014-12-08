unfollow.namespace 'unfollow.ajax', (exports) ->
  'use strict'

  call_function = (method, func, data, callback) ->
    if data
      data =
        data: JSON.stringify(data)

    request = $.ajax
      type: method
      url: '/ajax' + func
      dataType: 'json'
      data: data or {}
      beforeSend: (xhr) ->
        xhr.setRequestHeader 'X-CSRFToken', $.cookie('csrftoken') if method is 'POST'

    request.done (data) ->
      callback data, {} if callback
      location.href = data.redirect if data and data.redirect and data.redirect.slice(0, 1) isnt "/"
      $.pjax url: data.redirect if data and data.redirect

    request.error (xhr, txt_status) ->
      return if not callback
      callback null,
        ajax_error: txt_status

  exports.get = (func, data, callback) ->
    call_function 'GET', func, data, callback

  exports.post = (func, data, callback) ->
    call_function 'POST', func, data, callback

    #
    # ajax button
    #

  unfollow.init ->
    $('a.js-ajax, button.js-ajax').on 'click', ->
      button = $(this)
      ajax_get = button.data('ajax-get')
      ajax_post = button.data('ajax-post')
      return if not ajax_get and not ajax_post

      return if button.attr 'disabled'
      button.attr 'disabled', 'disabled'

      data = unfollow.utils.data button, (name) ->
        name == 'data-ajax-get' or name == 'data-ajax-post'

      callback = (data) ->
        return if data and (data.redirect or data.pjax)
        return button.remove() if data and data.remove
        button.removeAttr 'disabled'

      return unfollow.ajax.get ajax_get, data, callback if ajax_get
      return unfollow.ajax.post ajax_post, data, callback if ajax_post
  
