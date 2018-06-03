

$(function(){

    $("#about-btn").click( function(event) {
        alert("You clicked the button using JQuery!");
    });

});


$(function(){

    $(".list-group-item")
            .hover(
  function (event) {
    $(this).addClass('active');
  },
  function () {
    $(this).removeClass('active');
  }
  );

});