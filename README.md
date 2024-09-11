## Processador de Imagens em Go

Esse simples projeto visa demonstrar com um simples programa, como a utilização
de threads (sejam elas User Threads ou Kernel Threads) podem influenciar o tempo
de execução de uma dada tarefa. Especialmente se ela for [CPU-Bound](https://en.wikipedia.org/wiki/CPU-bound).

### Como funciona

Basicamente, o programa divide uma imagem em pedaços pequenos, onde podemos
aplicar um filtro e obter uma imagem de saída. Para fins de simplicidade, 
decidimos aplicar apenas uma operação de blurring (embaçamento).

Utilizamos as goroutines, que são unidades de execução gerenciadas pelo 
runtime do Go e multiplexadas em threads do sistema operacional, podendo
de certa foram ser consideradas como Threads de Usuário.

### Resultados

![Imagem Original](./assets/beach.jpg) ![Imagem Embaçada](./assets/beach_blurred.jpg)
