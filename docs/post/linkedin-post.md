# LinkedIn post draft — distributed-order-saga

Estou feliz em compartilhar meu projeto "distributed-order-saga": um sistema de pedidos distribuído que implementa o Saga Pattern (choreography-based) com Transactional Outbox. É um monorepo poliglota (Go + Ruby) e event-driven via Kafka.

Principais pontos:
- Arquitetura: Clean Architecture / DDD (domain/application/infrastructure)
- Resiliência: Transactional Outbox para publicar eventos de forma confiável
- Observabilidade: OpenTelemetry + Jaeger para traçar uma saga completa
- Serviços: order-service (Go), payment-service (Ruby), inventory-service (Go), notification-service (Ruby)

Demo e materiais:
- Código: https://github.com/joaodddev/distributed-order-saga
- Instruções para rodar localmente no README
- Screenshot do trace no Jaeger disponível no repositório

Se quiser ver um walkthrough em vídeo, me avise que eu gravo e adiciono ao repositório.

---

Sugestão de texto curto para o post:

"Lancei um projeto demonstrando o Saga Pattern com Transactional Outbox (Go + Ruby + Kafka). Arquitetura DDD, observabilidade com Jaeger e testes automatizados. Código e instruções em https://github.com/joaodddev/distributed-order-saga — feedbacks são bem-vindos!"
