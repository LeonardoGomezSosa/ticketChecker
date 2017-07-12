

//##############################< SCRIPTS JS >##########################################
//################################< Expresion.js >#####################################
//#########################< VALIDACIONES DE JEQUERY >##################################

$(document).ready(function () {
	console.log("-----------------");
	console.log("expresion.js");
	console.log("-----------------");
	$(".waitgif").hide();
	// var validator = valida();			
});

function valida() {
	var validator = $("#Form_Alta_Expresion").validate({
		rules: {

		},
		messages: {

		},
		errorElement: "em",
		errorPlacement: function (error, element) {
			error.addClass("help-block");
			element.parents(".col-sm-5").addClass("has-feedback");

			if (element.prop("type") === "checkbox") {
				error.insertAfter(element.parent("label"));
			} else {
				error.insertAfter(element);
			}

			if (!element.next("span")[0]) {
				$("<span class='glyphicon glyphicon-remove form-control-feedback'></span>").insertAfter(element);
			}
		},
		success: function (label, element) {
			if (!$(element).next("span")[0]) {
				$("<span class='glyphicon glyphicon-ok form-control-feedback'></span>").insertAfter($(element));
			}
		},
		highlight: function (element, errorClass, validClass) {
			$(element).parents(".col-sm-5").addClass("has-error").removeClass("has-success");
			$(element).next("span").addClass("glyphicon-remove").removeClass("glyphicon-ok");
		},
		unhighlight: function (element, errorClass, validClass) {
			$(element).parents(".col-sm-5").addClass("has-success").removeClass("has-error");
			$(element).next("span").addClass("glyphicon-ok").removeClass("glyphicon-remove");
		}
	});
	return validator;
}

function EditaExpresion(vista) {
	if (vista === "Index" || vista === "") {
		if ($('#Expresions').val() !== "") {
			window.location = '/Expresions/edita/' + $('#Expresions').val();
		} else {
			alertify.error("Debe Seleccionar un Expresion para editar");
		}
	} else if (vista == "Detalle") {
		console.log("editaExpresions: " + $('#IDExpresion').val());
		if ($('#IDExpresion').val() !== "") {
						alertify.confirm('Confirmar Edicion de Expresion regular', '¿Desea modificar el elemento ' + $('#IDExpresion').val() +
				' de Expresion Regular ' + $('#Expresion').val() + '?',
				function () {
					window.location = '/Expresions/edita/' + $('#IDExpresion').val();
				},
				function () {
					alertify.error('Accion Cancelada.');
				});
			// window.location = '/Expresions/edita/' + $('#IDExpresion').val();
		} else {
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Expresions';
		}
	}

}



function DetalleExpresion() {
	if ($('#IDExpresion').val() !== "") {
		window.location = '/Expresions/detalle/' + $('#IDExpresion').val();
	} else {
		alertify.error("Debe Seleccionar un Expresion para editar");
	}
}

function EliminaExpresion() {
	if ($('#IDExpresion').val() !== "") {
		alertify.confirm('Confirmar Eliminar', '¿Eliminar elemento ' + $('#IDExpresion').val() + ' de la expresion regular ' + $('#Expresion').val() + '?', function () { window.location = '/Expresions/Elimina/' + $('#IDExpresion').val(); }, function () { alertify.error('Accion Cancelada.') });
	} else {
		alertify.error("Debe Seleccionar un Expresion para Eliminar");
	}
}

function BuscaPagina(num) {
	$('#Loading').show();

	$.ajax({
		url: "/Expresions/search",
		type: 'POST',
		dataType: 'json',
		data: {
			Pag: num,
		},
		success: function (data) {
			if (data != null) {
				if (data.SEstado) {
					$("#Cabecera").empty();
					$("#Cabecera").append(data.SCabecera);
					$("#Cuerpo").empty();
					$("#Cuerpo").append(data.SBody);
					$("#Paginacion").empty();
					$("#Paginacion").append(data.SPaginacion);
				} else {
					alertify.error(data.SMsj);
				}
			} else {
				alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
			}
			$('#Loading').hide();
		},
		error: function (data) {
			$('#Loading').hide();
		},
	});
}


function SubmitGroup() {
	$('#Loading').show();
	$.ajax({
		url: "/Expresions/agrupa",
		type: 'POST',
		dataType: 'json',
		data: {
			Grupox: $('#Grupos').val(),
			searchbox: $('#searchbox').val()
		},
		success: function (data) {
			if (data != null) {
				if (data.SEstado) {
					$("#Cabecera").empty();
					$("#Cabecera").append(data.SCabecera);
					$("#Cuerpo").empty();
					$("#Cuerpo").append(data.SBody);
					$("#Paginacion").empty();
					$("#Paginacion").append(data.SPaginacion);
				} else {
					alertify.error(data.SMsj);
				}
			} else {
				alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
			}
			$('#Loading').hide();
		},
		error: function (data) {
			$('#Loading').hide();
		},
	});
}

