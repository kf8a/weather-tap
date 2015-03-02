$(function() {
  $('.graph').each(function(){
    var id = $(this).data('id')
    var label = $(this).data('label')
    var element = "#graph-"  + id

    console.log(element)
    var margin = {top: 20, right: 20, bottom: 30, left: 50},
      width = 960 - margin.left - margin.right,
      height = 200 - margin.top - margin.bottom;

    var parseTime = d3.time.format("%Y-%m-%dT%H:%M:%S").parse;

    var x = d3.time.scale()
      .range([0, width]);

    var y = d3.scale.linear()
      .range([height, 0]);

    var xAxis = d3.svg.axis()
      .scale(x)
      .orient("bottom");

    var yAxis = d3.svg.axis()
      .scale(y)
      .orient("left");

    var line = d3.svg.line()
      .x(function(d) { return x(d.time); })
      .y(function(d) { return y(d.value); });

    var svg = d3.select(this).append("svg")
      .attr("width", width + margin.left + margin.right)
      .attr("height", height + margin.top + margin.bottom)
      .append("g")
      .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    d3.json("/weather/variates/"+id, function(error, data) {
      data.forEach(function(d) {
        d.time = d.time.substring(0,19)
        d.time = parseTime(d.time);
        d.value = +d.value;
        console.log(d);
      });

      x.domain(d3.extent(data, function(d) { return d.time; }));
      y.domain(d3.extent(data, function(d) { return d.value; }));

      svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," + height + ")")
      .call(xAxis);

      svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)
      .append("text")
      .attr("transform", "rotate(-90)")
      .attr("y", 6)
      .attr("dy", ".71em")
      .style("text-anchor", "end")
      .text(label);

      svg.append("path")
      .datum(data)
      .attr("class", "line")
      .attr("d", line);
    })
  });
})
