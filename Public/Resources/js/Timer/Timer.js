$("document").ready(function () {
    $(".waitgif").hide();
    $('#RespuestaForm :input:enabled:visible:first').focus();
    $("#AceptarCodigo").on('click', function (e) {
        ValidarOperacion();
    });
    $("#AceptarRespuesta").on('click', function (e) {
        ValidarRespuesta();
    });
    $("#EntradaR").on('blur', function (event) {
        // if (event.keyCode == 13 || event.which == 13) {
            ValidarRespuesta();
        // }
    });
    $("#AceptarCodigo").on('focus', function (event) {
        // if (event.keyCode == 13 || event.which == 13) {
            ValidarOperacion();
        // }
    });

});

function ShowAlert() {
    $("#TimerDiv").alert();
}
function CloseAlert() {
    $("#TimerDiv").alert('close');
}
function ValidarRespuesta() {
    entrada = $("#EntradaR").val();
    $(".waitgif").show();
    var request = $.ajax({
        url: "/RecibirRespuesta",
        method: "POST",
        async: false,
        data: { Entrada: entrada },
        dataType: "html",
    });
    request.done(function (data) {
        $("body").html(data);
    });
    request.fail(function (data) {
        $("body").html(data);
    });
    request.always(function () {
        $("#EntradaR").val("");
        timerOn = $("#TimerOn").val();
        console.log(timerOn);
        $(".waitgif").hide();
        if (timerOn === true) {
            setTimeout(CloseAlert, 3000);
        }
        $('#RespuestaForm :input:enabled:visible:first').focus();
    });
}

function ValidarOperacion() {
    entrada = $("#Entrada").val();
    ticket = $("#Ticket").val();
    surtidor = $("#Surtidor").val();
    timer = $("#TimerOn").val();
    $(".waitgif").show();

    var request = $.ajax({
        url: "/",
        method: "POST",
        async: false,
        data: { Entrada: entrada, Ticket: ticket, Surtidor: surtidor , TimerOn:timer},
        dataType: "html",
    });
    request.done(function (data) {
        $("body").html(data);
    });
    request.fail(function (data) {
        $("body").html(data);
    });
    request.always(function () {
        $("#Entrada").val("");
        ticket = $("#Ticket").val();
        surtidor = $("#Surtidor").val();
        timerOn = $("#TimerOn").val();

        setTimeout(CloseAlert, 3000);
        $('#RespuestaForm :input:enabled:visible:first').focus();
    });
}
