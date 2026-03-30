# Objetivo

Quiero crear una aplicacion de CLI con Golang y debe tener una TUI, habia pensado en bubbletea,

El objetivo es tener una aplicacion que pueda ejecutar y me ayude a configurar mi sistema operativo instalando paquetes y aplicando configuracionies. Paso a paso.

Queremos dar soportes a diferentes sistemas operativos pero todos basados en Linux. Comenzaremos unicamente con Debian, ten en cuenta que puede haber diferencias de paquetes entre versiones, de modo que tambien hay que distinguir entre Debian 12 o 13 en este momento), de modo que lo primero será crear la estructura necesaria para identificar el sistema y comenzar el proceso. Tambien debido a que vamos a instalar paquetes hay que comprobar si se está ejecutando con sudo o como usuario root. No se debe permitir su ejecución sin sudo o root.

Vamos a dividir el proceso en diferentes pasos y el usuario irá confirmando y escogiendo opciones.

Para empezar, vamos a dar soporte para instalar navegadores, será el primer paso donde el usuario pueda interactuar. Hay que dar la opcion para instalar Firefox, Google Chrome, Chromium y Brave, el usuario podrá escoger cuales quiere con opcion multiple y habrá que instalarlos desde la instalación oficial.

Por ejemplo para Brave he encontado este codigo:
"""
sudo apt install curl

sudo curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg

sudo curl -fsSLo /etc/apt/sources.list.d/brave-browser-release.sources https://brave-browser-apt-release.s3.brave.com/brave-browser.sources

sudo apt update

sudo apt install brave-browser
"""

Localiza ellas opciones para el resto de navegadores.

Absolutamente todo debe tenes testing, y hay que desarrollar absolutamente todo con TDD.



