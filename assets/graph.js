$(function() {
  $('.graph').each(function(){
    var id = $(this).data('id')
    var element = "#graph-"  + id
    console.log(element)
    console.log(c3)
    console.log(d3)
    var chart = c3.generate({
      bindto: element,
      data: {
        columns: [
          ['data1', 30, 200, 100, 400, 150, 250],
          ['data2', 50, 20, 10, 40, 15, 25]
        ]
      }

      // data: {
      //   url: "variates/"+id,
      //   type: 'line'
      // }
    });
  });
})
