$(document).ready(function () {

		$('#inputBuscaBase').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();							
				ConsultaBase();
		}
        });

        $('#inputBuscaSat').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();							
				ConsultaSat();
		}
        });


    $('input[type=radio][name=filtroBase]').change(function(e) {
		SubmitGroupBase();
    });

});

$(document).ready(function() {
  $(window).keydown(function(event){
    if(event.keyCode == 13) {
      event.preventDefault();
    }
  });
});

//Funciones para el manejo de los checkbox en las vistas
//var SkuSeleccionados = [];
var SkuSeleccionados = new Array();
var ClaveSatSeleccionado;

function GuardaSeleccionados(seleccionado, valor){
    if (seleccionado==true){
        SkuSeleccionados.push(valor);
    }else{
        var indiceEncontrado = SkuSeleccionados.indexOf(valor);
        removeA(SkuSeleccionados, valor);
    }	
}

/**
 * Funcion que seleciona o deselecciona  los checbox de productos
 * @param {*} seleccionado  El checbox de todos
 */
function SelecionarTodos(seleccionado){

	if (seleccionado==true){
        $( ".ProdExtraido" ).each(function() {
			var id=$(this).attr( "id" );
  			$(this).prop('checked', true);
			  SkuSeleccionados.push(id);
		});
    }else{
		$( ".ProdExtraido" ).each(function() {
			var id=$(this).attr( "id");
  			$(this).prop('checked', false);
			  removeA(SkuSeleccionados, id);
		});
    }	
}


function ClickSelecionarTodos(seleccionado){
	if (seleccionado==true){
        SkuSeleccionados.push(valor);
    }else{
        var indiceEncontrado = SkuSeleccionados.indexOf(valor);
        removeA(SkuSeleccionados, valor);
    }	
}

function Verificarseleccionados(){
    for (var x=0; x<SkuSeleccionados.length; x++){
		if(document.getElementById(SkuSeleccionados[x])) {
			document.getElementById(SkuSeleccionados[x]).setAttribute("checked", "checked");
		}
    }
	var contador=1;
	$( ".ProdExtraido" ).each(function() {
		if( !$(this).is(':checked') ) {
    		contador=0;
			return false;
		}
	});

	if(contador==1){
		document.getElementById("selectTodos").setAttribute("checked", "checked");
	}
	
}

function removeA(arr) {
    var what, a = arguments, L = a.length, ax;
    while (L > 1 && arr.length) {
        what = a[--L];
        while ((ax= arr.indexOf(what)) !== -1) {
            arr.splice(ax, 1);
        }
    }
    return arr;
}

function AsignarClaveSat(valor){
	ClaveSatSeleccionado = valor;
}

function AgregarSkuSeleccionadosAEnviar(){
	if (SkuSeleccionados.length>0){
		$("#SkuSeleccionados").val(JSON.stringify(SkuSeleccionados));
		return true;
	}else{
		return false;
	}
}

function AgregarClaveSeleccionadosAEnviar(){
	if (ClaveSatSeleccionado!=""){
		$("#ClaveSatSeleccionado").val(ClaveSatSeleccionado);
		console.log(ClaveSatSeleccionado);
		return true;
	}else{
		return false;
	}	
}
//Fin de los metodos de checbox

function ConsultaBase(){
    var cadena = $("#inputBuscaBase").val();
	SkuSeleccionados=[];
    if (cadena != ""){
        	 $('#Loading').show();

        $.ajax({
			url:"/Listas/ConsultaBase",
			type: 'POST',
			dataType:'json',
			data:{
				Cadena : cadena,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBase").empty();		
						$("#PaginacionBase").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
	 				$('#Loading').hide(); 
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
            				$('#Loading').hide(); 
		  },
		});
    }else{
        alertify.error("Introduce una cadena de texto a consultar.");
        $("#inputBuscaBase").focus();
    }
}

function ConsultaSat(){
    var cadena = $("#inputBuscaSat").val();

    if (cadena != ""){
        	 $('#Loading').show();
			 ClaveSatSeleccionado="";
        $.ajax({
			url:"/Listas/ConsultaSat",
			type: 'POST',
			dataType:'json',
			data:{
				Cadena : cadena,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraSat").empty();						
						$("#CabeceraSat").append(data.SCabecera);
						$("#BodySat").empty();						
						$("#BodySat").append(data.SBody);
						$("#PaginacionSat").empty();		
						$("#PaginacionSat").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
	 				$('#Loading').hide(); 
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
            				$('#Loading').hide(); 
		  },
		});
    }else{
        alertify.error("Introduce una cadena de texto a consultar.");
        $("#inputBuscaBase").focus();
    }
}



function BuscaPagina(num){
			$.ajax({
			url:"/Listas/search",
			type: 'POST',
			dataType:'json',
			data:{
				Pag : num,
				Filtro: $('input:checked[type=radio][name=filtroBase]').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBase").empty();		
						$("#PaginacionBase").append(data.SPaginacion);
						Verificarseleccionados();				
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
	 
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
		  },
		});
}

function BuscaPaginaSat(num){
			$.ajax({
			url:"/Listas/searchsat",
			type: 'POST',
			dataType:'json',
			data:{
				Pag : num,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraSat").empty();						
						$("#CabeceraSat").append(data.SCabecera);
						$("#BodySat").empty();						
						$("#BodySat").append(data.SBody);
						$("#PaginacionSat").empty();		
						$("#PaginacionSat").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
	 
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
		  },
		});
}

 function SubmitGroupBase(){
	 $('#Loading').show();
			$.ajax({
			url:"/Listas/agrupaB",
			type: 'POST',
			dataType:'json',
			data:{
				GrupoBase : $('#GruposBase').val(),
				Cadena: $('#inputBuscaBase').val(),
				Filtro: $('input:checked[type=radio][name=filtroBase]').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBase").empty();		
						$("#PaginacionBase").append(data.SPaginacion);					
					}else{
						$("#CabeceraBase").empty();						
						$("#BodyBase").empty();						
						$("#PaginacionBase").empty();				
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				$('#Loading').hide(); 
			},
		  error: function(data){
			  $('#Loading').hide();
		  },
		});
}

 function SubmitGroupSat(){
	 $('#Loading').show();
			$.ajax({
			url:"/Listas/agrupaS",
			type: 'POST',
			dataType:'json',
			data:{
				GrupoSat : $('#GruposSat').val(),
				Cadena: $('#inputBuscaSat').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraSat").empty();						
						$("#CabeceraSat").append(data.SCabecera);
						$("#BodySat").empty();						
						$("#BodySat").append(data.SBody);
						$("#PaginacionSat").empty();		
						$("#PaginacionSat").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				$('#Loading').hide(); 
			},
		  error: function(data){
			  $('#Loading').hide();
		  },
		});
}
/**
 * Metodo de descarga de catalogo a archivo CSV
 */
function descargaCatalogo(){

  alertify.confirm("Información","¿Confirmar que desea descargar el catalogo?",
  function(){
		$.ajax({
			url: '/Descarga',
			type: 'GET',
			dataType: 'html',
			data:{},
			success : function(data){
    			alertify.success('Ok');
				$("#contenido").html(data);
				//$("#downloadlink").html("<a href='http://localhost:9095/Download/data.csv' download='catalogo'>Descargar Archivo</a>")				
			}
		});
  },
  function(){
    alertify.error('Cancel');
  });
  /*	
	*/	
}
