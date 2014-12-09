unfollow.namespace 'unfollow.utils', (exports) ->
  'use strict'

  # setup google analytics
  analytics = (account) ->
    return unless account

    require ['ga'], ->
      window.ga 'create', account, 'unfollow.io'
      window.ga 'send', 'pageview',
        page: location.pathname + location.search + location.hash
        location: location.href

  # extract data attributes on an element into an object
  exports.data = (element, filter) ->
    filter = filter or -> false
    data = {}
    $.each element[0].attributes, (index, attribute) ->
      return unless /^data-/.test(attribute.name)
      return if filter(attribute.name)
      data[attribute.name.substr(5).replace(/\-/g, '_')] = attribute.value
    return data

  exports.show_modal = (modal, callback) ->
    $('.modal:visible').modal('hide')

    modal = $(modal).removeClass('hide').addClass('offscreen').show().appendTo $(document.body)
    setTimeout ->
      modal.removeClass('offscreen').modal()
      modal.on 'hidden.bs.modal', ->
        modal.remove()
      modal.on 'shown.bs.modal', ->
        modal.find('input:visible,textarea:visible').filter(':first').focus()
      callback modal if callback
    , 150

  unfollow.init ->
    analytics(unfollow.settings.ANALYTICS)

    $('div.js-modal').livequery ->
      modal = $(this).removeClass('hide')
      modal.modal().on 'hidden.bs.modal', -> modal.remove()
