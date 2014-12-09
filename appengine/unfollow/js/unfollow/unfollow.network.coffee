unfollow.namespace 'unfollow.network', (exports) ->
  'use strict'

  class Network
    constructor: (@network) ->
      @nodes = [
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {}
      ]
      @links = [
        {source:1, target:0},
        {source:2, target:0},
        {source:3, target:0},
        {source:3, target:2},
        {source:4, target:0},
        {source:5, target:0},
        {source:6, target:0},
        {source:7, target:0},
        {source:8, target:0},
        {source:9, target:0},
        {source:11, target:10},
        {source:11, target:3},
        {source:11, target:2},
        {source:11, target:0},
        {source:12, target:11},
        {source:13, target:11},
        {source:14, target:11},
        {source:15, target:11},
        {source:17, target:16},
        {source:18, target:16},
        {source:18, target:17},
        {source:19, target:16},
        {source:19, target:17},
        {source:19, target:18},
        {source:20, target:16},
        {source:20, target:17},
        {source:20, target:18},
        {source:20, target:19},
        {source:21, target:16},
        {source:21, target:17},
        {source:21, target:18},
        {source:21, target:19},
        {source:21, target:20},
        {source:22, target:16},
        {source:22, target:17},
        {source:22, target:18},
        {source:22, target:19},
        {source:22, target:20},
        {source:22, target:21},
        {source:23, target:16},
        {source:23, target:17},
        {source:23, target:18},
        {source:23, target:19},
        {source:23, target:20},
        {source:23, target:21},
        {source:23, target:22},
        {source:23, target:12},
        {source:23, target:11},
        {source:24, target:23},
        {source:24, target:11},
        {source:25, target:24},
        {source:25, target:23},
        {source:25, target:11},
        {source:26, target:24},
        {source:26, target:11},
        {source:26, target:16},
        {source:26, target:25},
        {source:27, target:11},
        {source:27, target:23},
        {source:27, target:25},
        {source:27, target:24},
        {source:27, target:26},
        {source:28, target:11},
        {source:28, target:27},
        {source:29, target:23},
        {source:29, target:27},
        {source:29, target:11},
        {source:30, target:23},
        {source:31, target:30},
        {source:31, target:11},
        {source:31, target:23},
        {source:31, target:27},
        {source:32, target:11},
        {source:33, target:11},
        {source:33, target:27},
        {source:34, target:11},
        {source:34, target:29},
        {source:35, target:11},
        {source:35, target:34},
        {source:35, target:29},
        {source:36, target:34},
        {source:36, target:35},
        {source:36, target:11},
        {source:36, target:29},
        {source:37, target:34},
        {source:37, target:35},
        {source:37, target:36},
        {source:37, target:11},
        {source:37, target:29},
        {source:38, target:34},
        {source:38, target:35},
        {source:38, target:36},
        {source:38, target:37},
        {source:38, target:11},
        {source:38, target:29},
        {source:39, target:25},
        {source:40, target:25},
        {source:41, target:24},
        {source:41, target:25},
        {source:42, target:41},
        {source:42, target:25},
        {source:42, target:24},
        {source:43, target:11},
        {source:43, target:26},
        {source:43, target:27},
        {source:44, target:28},
        {source:44, target:11},
        {source:45, target:28},
        {source:47, target:46},
        {source:48, target:47},
        {source:48, target:25},
        {source:48, target:27},
        {source:48, target:11},
        {source:49, target:26},
        {source:49, target:11},
        {source:50, target:49},
        {source:50, target:24},
        {source:51, target:49},
        {source:51, target:26},
        {source:51, target:11},
        {source:52, target:51},
        {source:52, target:39},
        {source:53, target:51},
        {source:54, target:51},
        {source:54, target:49},
        {source:54, target:26},
        {source:55, target:51},
        {source:55, target:49},
        {source:55, target:39},
        {source:55, target:54},
        {source:55, target:26},
        {source:55, target:11},
        {source:55, target:16},
        {source:55, target:25},
        {source:55, target:41},
        {source:55, target:48},
        {source:56, target:49},
        {source:56, target:55},
        {source:57, target:55},
        {source:57, target:41},
        {source:57, target:48},
        {source:58, target:55},
        {source:58, target:48},
        {source:58, target:27},
        {source:58, target:57},
        {source:58, target:11},
        {source:59, target:58},
        {source:59, target:55},
        {source:59, target:48},
        {source:59, target:57},
        {source:60, target:48},
        {source:60, target:58},
        {source:60, target:59},
        {source:61, target:48},
        {source:61, target:58},
        {source:61, target:60},
        {source:61, target:59},
        {source:61, target:57},
        {source:61, target:55},
        {source:62, target:55},
        {source:62, target:58},
        {source:62, target:59},
        {source:62, target:48},
        {source:62, target:57},
        {source:62, target:41},
        {source:62, target:61},
        {source:62, target:60},
        {source:63, target:59},
        {source:63, target:48},
        {source:63, target:62},
        {source:63, target:57},
        {source:63, target:58},
        {source:63, target:61},
        {source:63, target:60},
        {source:63, target:55},
        {source:64, target:55},
        {source:64, target:62},
        {source:64, target:48},
        {source:64, target:63},
        {source:64, target:58},
        {source:64, target:61},
        {source:64, target:60},
        {source:64, target:59},
        {source:64, target:57},
        {source:64, target:11},
        {source:65, target:63},
        {source:65, target:64},
        {source:65, target:48},
        {source:65, target:62},
        {source:65, target:58},
        {source:65, target:61},
        {source:65, target:60},
        {source:65, target:59},
        {source:65, target:57},
        {source:65, target:55},
        {source:66, target:64},
        {source:66, target:58},
        {source:66, target:59},
        {source:66, target:62},
        {source:66, target:65},
        {source:66, target:48},
        {source:66, target:63},
        {source:66, target:61},
        {source:66, target:60},
        {source:67, target:57},
        {source:68, target:25},
        {source:68, target:11},
        {source:68, target:24},
        {source:68, target:27},
        {source:68, target:48},
        {source:68, target:41},
        {source:69, target:25},
        {source:69, target:68},
        {source:69, target:11},
        {source:69, target:24},
        {source:69, target:27},
        {source:69, target:48},
        {source:69, target:41},
        {source:70, target:25},
        {source:70, target:69},
        {source:70, target:68},
        {source:70, target:11},
        {source:70, target:24},
        {source:70, target:27},
        {source:70, target:41},
        {source:70, target:58},
        {source:71, target:27},
        {source:71, target:69},
        {source:71, target:68},
        {source:71, target:70},
        {source:71, target:11},
        {source:71, target:48},
        {source:71, target:41},
        {source:71, target:25},
        {source:72, target:26},
        {source:72, target:27},
        {source:72, target:11},
        {source:73, target:48},
        {source:74, target:48},
        {source:74, target:73},
        {source:75, target:69},
        {source:75, target:68},
        {source:75, target:25},
        {source:75, target:48},
        {source:75, target:41},
        {source:75, target:70},
        {source:75, target:71},
        {source:76, target:64},
        {source:76, target:65},
        {source:76, target:66},
        {source:76, target:63},
        {source:76, target:62},
        {source:76, target:48},
        {source:76, target:58}
      ]

      @setup()

    setup: ->
      @svg = d3.select(@network[0]).append('svg')
      @defs = @svg.append('defs')
      @overlay = @svg.append('rect').attr('class', 'overlay')
      @g = @svg.append('g')

      @link = @g.selectAll('.link')
      @node = @g.selectAll('.node')
      @svg.call d3.behavior.zoom().scaleExtent([1, 10]).on 'zoom', =>
        @g.attr('transform', 'translate(' + d3.event.translate + ')scale(' + d3.event.scale + ')')

      @force = d3.layout.force().charge(-120).linkDistance(30)
      @force.on 'tick', =>
        @link
          .attr('x1', (data) -> data.source.x)
          .attr('y1', (data) -> data.source.y)
          .attr('x2', (data) -> data.target.x)
          .attr('y2', (data) -> data.target.y)
        @node
          .attr('cx', (data) -> data.x)
          .attr('cy', (data) -> data.y)
      @force.drag()
        .on 'dragstart', ->
          d3.event.sourceEvent.stopPropagation()
        .on 'dragend', (data) ->
          data.fixed = true

      $(window).resize => @resize()
      @resize()

      @defs
        .append('pattern')
        .attr('id', 'ziyan')
        .attr('patternUnits', 'objectBoundingBox')
        .attr('width', 10)
        .attr('height', 10)
        .append('image')
        .attr('xlink:href', 'https://pbs.twimg.com/profile_images/1169934089/avatarpic-l_normal.png')
        .attr('x', 0)
        .attr('y', 0)
        .attr('width', 10)
        .attr('height', 10)

    resize: ->
      width = $(window).width()
      height = $(window).height()

      @svg.attr('width', width).attr('height', height)
      @overlay.attr('width', width).attr('height', height)
      @force.size([width, height])

      @update()

    update: ->
      @force.nodes(@nodes).links(@links).start()

      @link = @link.data(@links)
      @link.enter().append('line')
        .attr('class', 'link')
        .style('stroke-width', 1)
      @link.exit().remove()

      @node = @node.data(@nodes)
      @node.enter().append('circle')
        .attr('class', 'node')
        .attr('r', 5)
        .attr('style', 'fill: url(#ziyan)')
        .call(@force.drag)
      @node.exit().remove()

  on_network_init = ->
    network = $(this)
    network.data 'network', new Network(network)

  unfollow.init ->
    $('div.js-network').each on_network_init
