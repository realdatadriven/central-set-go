<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# PAREP_MODEL_V2
```yaml
name: PAREPCV
description: Modelo de gestão do programa PAREP-CV para objetivos, atividades, indicadores, ações, contratos, entidades e parametrização.
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
depends_on: ADMIN
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  PAREPCV:
    menu_icon: chart-pie
    menu_order: 1
    active: true
    tables:
      - resultados
      - atividades
      - atividade_anos
      - acoes
      - contratos
      - contrato_execucoes
      - contrato_anexos
      - logistica
      - beneficiarios
      - indicadores
      - metas
      - entidades
      - anos
      - ilhas
      - concelhos
      - generos
      - tipo_contrato
      - tipo_entidade
      - status_acao
      - status_contrato
      - status_execucao
      - unidades
  Parametrização:
    menu_icon: adjustments
    menu_order: 2
    active: true
    tables:
      - anos
      - ilhas
      - concelhos
      - generos
      - tipo_contrato
      - tipo_entidade
      - status_acao
      - status_contrato
      - status_execucao
      - unidades
```

## RESULTADOS
```yaml
table: resultados
comment: Resultado
tooltip: Resultado estratégico do programa PAREP-CV e respetiva orientação de aprendizagem.
columns:
  resultado_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador único do resultado.", form_display: true, table_display: true, order: 1 }
  resultado: { type: varchar, len: 255, nullable: false, comment: "Resultado", tooltip: "Nome do resultado estratégico.", form_display: true, table_display: true, form_size: 8, order: 2 }
  desc_resultado: { type: text, comment: "Descrição", tooltip: "Descrição detalhada do resultado.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início do período de execução do resultado.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de término do período de execução do resultado.", form_display: true, table_display: true, form_size: 4, order: 5 }
  status_id: { type: integer, comment: "Estado", tooltip: "Estado do resultado.", form_display: true, table_display: true, form_size: 4, order: 6 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {atividades: true, indicadores: true}
table_layout:
  default_order: [{field: resultado_id, order: DESC}]
```

## ATIVIDADES
```yaml
table: atividades
comment: Atividade
tooltip: Atividade associada a um resultado e ao programa de intervenção.
columns:
  atividade_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da atividade.", form_display: true, table_display: true, order: 1 }
  atividade: { type: varchar, len: 255, nullable: false, comment: "Atividade", tooltip: "Nome da atividade.", form_display: true, table_display: true, form_size: 8, order: 2 }
  desc_atividade: { type: text, comment: "Descrição", tooltip: "Descrição da atividade.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", nullable: false, comment: "Resultado", tooltip: "Resultado ao qual a atividade está relacionada.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da atividade.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da atividade.", form_display: true, table_display: true, form_size: 4, order: 6 }
  status_id: { type: integer, comment: "Estado", tooltip: "Estado da atividade.", form_display: true, table_display: true, form_size: 4, order: 7 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {atividade_anos: true}
table_layout:
  default_order: [{field: atividade_id, order: DESC}]
```

## ATIVIDADE_ANOS
```yaml
table: atividade_anos
comment: Atividade Ano
tooltip: Variação anual de uma atividade associada ao ano e respetiva execução.
columns:
  atividade_ano_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do registo anual da atividade.", form_display: true, table_display: true, order: 1 }
  atividade_id: { type: integer, fk: "atividades.atividade_id", nullable: false, comment: "Atividade", tooltip: "Atividade principal.", form_display: true, table_display: true, form_size: 4, order: 2 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", nullable: false, comment: "Resultado", tooltip: "Resultado associado.", form_display: true, table_display: true, form_size: 4, order: 3 }
  ano_id: { type: integer, fk: "anos.ano_id", nullable: false, comment: "Ano", tooltip: "Ano de execução da atividade.", form_display: true, table_display: true, form_size: 4, order: 4 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano literal para referência rápida.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da execução anual.", form_display: true, table_display: true, form_size: 4, order: 6 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da execução anual.", form_display: true, table_display: true, form_size: 4, order: 7 }
  meta_ano: { type: decimal, len: 14, scale: 2, default: 0, comment: "Meta anual", tooltip: "Meta de execução anual da atividade.", form_display: true, table_display: true, form_size: 4, order: 8 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {acoes: true}
table_layout:
  default_order: [{field: atividade_ano_id, order: DESC}]
```

## ACOES
```yaml
table: acoes
comment: Ação
tooltip: Ação de implementação associada a uma atividade anual.
columns:
  acao_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da ação.", form_display: true, table_display: true, order: 1 }
  acao: { type: varchar, len: 255, nullable: false, comment: "Ação", tooltip: "Nome da ação.", form_display: true, table_display: true, form_size: 8, order: 2 }
  acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição da ação.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da ação.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da ação.", form_display: true, table_display: true, form_size: 4, order: 5 }
  orcamento_estimado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento estimado", tooltip: "Valor estimado da ação.", form_display: true, table_display: true, form_size: 4, order: 6 }
  orcamento_executado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento executado", tooltip: "Valor executado da ação.", form_display: true, table_display: true, form_size: 4, order: 7 }
  status_id: { type: integer, fk: "status_acao.status_acao_id", comment: "Status", tooltip: "Estado atual da ação.", form_display: true, table_display: true, form_size: 4, order: 8 }
  atividade_ano_id: { type: integer, fk: "atividade_anos.atividade_ano_id", comment: "Atividade ano", tooltip: "Relação com a execução anual da atividade.", form_display: true, table_display: true, form_size: 4, order: 9 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano da ação for referência.", form_display: true, table_display: true, form_size: 4, order: 10 }
  atividade_id: { type: integer, fk: "atividades.atividade_id", comment: "Atividade", tooltip: "Atividade associada.", form_display: true, table_display: true, form_size: 4, order: 11 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado", tooltip: "Resultado associado.", form_display: true, table_display: true, form_size: 4, order: 12 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {contratos: true, logistica: true, beneficiarios: true}
table_layout:
  default_order: [{field: acao_id, order: DESC}]
```

## CONTRATOS
```yaml
table: contratos
comment: Contrato
tooltip: Contrato associado a uma ação e ao respetivo acompanhamento de execução.
columns:
  contrato_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do contrato.", form_display: true, table_display: true, order: 1 }
  contrato: { type: varchar, len: 255, nullable: false, comment: "Contrato", tooltip: "Nome do contrato.", form_display: true, table_display: true, form_size: 8, order: 2 }
  codigo_contrato: { type: varchar, len: 100, comment: "Código", tooltip: "Código ou referência do contrato.", form_display: true, table_display: true, form_size: 4, order: 3 }
  desc_contrato: { type: text, comment: "Descrição", tooltip: "Descrição do contrato.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  data_ini_contrato: { type: date, comment: "Data início", tooltip: "Data de início do contrato.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim_contrato: { type: date, comment: "Data fim", tooltip: "Data de fim do contrato.", form_display: true, table_display: true, form_size: 4, order: 6 }
  entidade_id: { type: integer, fk: "entidades.entidade_id", comment: "Entidade", tooltip: "Entidade contratada.", form_display: true, table_display: true, form_size: 4, order: 7 }
  acao_id: { type: integer, fk: "acoes.acao_id", nullable: false, comment: "Ação", tooltip: "Ação a que o contrato está associado.", form_display: true, table_display: true, form_size: 4, order: 8 }
  atividade_ano_id: { type: integer, fk: "atividade_anos.atividade_ano_id", comment: "Atividade ano", tooltip: "Relação com a atividade anual.", form_display: true, table_display: true, form_size: 4, order: 9 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano do contrato.", form_display: true, table_display: true, form_size: 4, order: 10 }
  atividade_id: { type: integer, fk: "atividades.atividade_id", comment: "Atividade", tooltip: "Atividade associada.", form_display: true, table_display: true, form_size: 4, order: 11 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado", tooltip: "Resultado associado.", form_display: true, table_display: true, form_size: 4, order: 12 }
  status_id: { type: integer, fk: "status_contrato.status_contrato_id", comment: "Status", tooltip: "Estado do contrato.", form_display: true, table_display: true, form_size: 4, order: 13 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {contrato_execucoes: true, contrato_anexos: true}
table_layout:
  default_order: [{field: contrato_id, order: DESC}]
```

## CONTRATO_EXECUCOES
```yaml
table: contrato_execucoes
comment: Execução do Contrato
tooltip: Fases ou execuções do contrato com montante e acompanhamento temporal.
columns:
  contrato_execucao_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da execução do contrato.", form_display: true, table_display: true, order: 1 }
  contrato_execucao: { type: varchar, len: 255, comment: "Execução", tooltip: "Nome da fase ou execução do contrato.", form_display: true, table_display: true, form_size: 6, order: 2 }
  desc_contrato_execucao: { type: text, comment: "Descrição", tooltip: "Descrição da execução do contrato.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  montante_executado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Montante executado", tooltip: "Montante já executado pela fase.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_ini_fase: { type: date, comment: "Data início fase", tooltip: "Data de início da fase.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim_fase: { type: date, comment: "Data fim fase", tooltip: "Data de fim da fase.", form_display: true, table_display: true, form_size: 4, order: 6 }
  status_execucao_id: { type: integer, fk: "status_execucao.status_execucao_id", comment: "Status execução", tooltip: "Estado da fase de execução.", form_display: true, table_display: true, form_size: 4, order: 7 }
  attach_relatorio_contrato: { type: varchar, len: 255, comment: "Relatório", tooltip: "Anexo do relatório de execução do contrato.", form_display: true, table_display: true, form_size: 6, order: 8 }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato associado.", form_display: true, table_display: true, form_size: 6, order: 9 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
table_layout:
  default_order: [{field: contrato_execucao_id, order: DESC}]
```

## CONTRATO_ANEXOS
```yaml
table: contrato_anexos
comment: Anexo do Contrato
tooltip: Documentos ou anexos associados ao contrato.
columns:
  contrato_anexo_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do anexo do contrato.", form_display: true, table_display: true, order: 1 }
  contrato_anexo: { type: varchar, len: 255, comment: "Anexo", tooltip: "Nome do anexo.", form_display: true, table_display: true, form_size: 6, order: 2 }
  attach_contrato_anexo: { type: varchar, len: 255, comment: "Ficheiro", tooltip: "Nome ou caminho do ficheiro anexado.", form_display: true, table_display: true, form_size: 6, order: 3 }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato ao qual o anexo pertence.", form_display: true, table_display: true, form_size: 6, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order: [{field: contrato_anexo_id, order: DESC}]
```

## LOGISTICA
```yaml
table: logistica
comment: Logística
tooltip: Gestão logística associada à atividade anual e respetivo acompanhamento de custos.
columns:
  logistica_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da logística.", form_display: true, table_display: true, order: 1 }
  logistica: { type: varchar, len: 255, nullable: false, comment: "Logística", tooltip: "Descrição da logística ou operação.", form_display: true, table_display: true, form_size: 8, order: 2 }
  orcamento_estimado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento estimado", tooltip: "Estimativa de custos da logística.", form_display: true, table_display: true, form_size: 4, order: 3 }
  orcamento_executado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento executado", tooltip: "Valor executado da logística.", form_display: true, table_display: true, form_size: 4, order: 4 }
  status_id: { type: integer, fk: "status_acao.status_acao_id", comment: "Status", tooltip: "Estado da logística.", form_display: true, table_display: true, form_size: 4, order: 5 }
  attach_logistica_anexo: { type: varchar, len: 255, comment: "Anexo", tooltip: "Anexo ou documento relacionado com a logística.", form_display: true, table_display: true, form_size: 6, order: 6 }
  atividade_ano_id: { type: integer, fk: "atividade_anos.atividade_ano_id", comment: "Atividade ano", tooltip: "Atividade anual associada.", form_display: true, table_display: true, form_size: 6, order: 7 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano da logística.", form_display: true, table_display: true, form_size: 4, order: 8 }
  atividade_id: { type: integer, fk: "atividades.atividade_id", comment: "Atividade", tooltip: "Atividade principal.", form_display: true, table_display: true, form_size: 4, order: 9 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado", tooltip: "Resultado associado.", form_display: true, table_display: true, form_size: 4, order: 10 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: logistica_id, order: DESC}]
```

## BENEFICIARIOS
```yaml
table: beneficiarios
comment: Beneficiário
tooltip: Registo dos beneficiários associados a uma atividade anual.
columns:
  beneficiario_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do beneficiário.", form_display: true, table_display: true, order: 1 }
  beneficiario: { type: varchar, len: 255, nullable: false, comment: "Beneficiário", tooltip: "Nome do beneficiário.", form_display: true, table_display: true, form_size: 8, order: 2 }
  data_ini_benef: { type: date, comment: "Data início", tooltip: "Data de início da elegibilidade ou apoio.", form_display: true, table_display: true, form_size: 4, order: 3 }
  data_fim_benef: { type: date, comment: "Data fim", tooltip: "Data de fim do apoio ou acompanhamento.", form_display: true, table_display: true, form_size: 4, order: 4 }
  beneficiario_status_id: { type: integer, comment: "Status", tooltip: "Status do beneficiário.", form_display: true, table_display: true, form_size: 4, order: 5 }
  attach_beneficiarios: { type: varchar, len: 255, comment: "Anexo", tooltip: "Anexo ou ficheiro associado ao beneficiário.", form_display: true, table_display: true, form_size: 6, order: 6 }
  atividade_ano_id: { type: integer, fk: "atividade_anos.atividade_ano_id", comment: "Atividade ano", tooltip: "Atividade anual associada.", form_display: true, table_display: true, form_size: 6, order: 7 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano de referência.", form_display: true, table_display: true, form_size: 4, order: 8 }
  atividade_id: { type: integer, fk: "atividades.atividade_id", comment: "Atividade", tooltip: "Atividade principal.", form_display: true, table_display: true, form_size: 4, order: 9 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado", tooltip: "Resultado associado.", form_display: true, table_display: true, form_size: 4, order: 10 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: beneficiario_id, order: DESC}]
```

## INDICADORES
```yaml
table: indicadores
comment: Indicador
tooltip: Indicador de resultado ou de acompanhamento do desempenho do programa.
columns:
  indicador_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do indicador.", form_display: true, table_display: true, order: 1 }
  indicador: { type: varchar, len: 255, nullable: false, comment: "Indicador", tooltip: "Nome do indicador.", form_display: true, table_display: true, form_size: 8, order: 2 }
  desc_indicador: { type: text, comment: "Descrição", tooltip: "Descrição detalhada do indicador.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  ano_baseline: { type: integer, comment: "Ano baseline", tooltip: "Ano de referência para a linha de base.", form_display: true, table_display: true, form_size: 4, order: 4 }
  valor_baseline: { type: decimal, len: 14, scale: 2, default: 0, comment: "Valor baseline", tooltip: "Valor inicial de referência.", form_display: true, table_display: true, form_size: 4, order: 5 }
  unidades_id: { type: integer, fk: "unidades.unidade_id", comment: "Unidade", tooltip: "Unidade de medida do indicador.", form_display: true, table_display: true, form_size: 4, order: 6 }
  fonte: { type: varchar, len: 255, comment: "Fonte", tooltip: "Fonte ou origem dos dados do indicador.", form_display: true, table_display: true, form_size: 6, order: 7 }
  formula_calculo: { type: text, comment: "Fórmula", tooltip: "Fórmula de cálculo do indicador.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 8 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", nullable: false, comment: "Resultado", tooltip: "Resultado ao qual o indicador pertence.", form_display: true, table_display: true, form_size: 6, order: 9 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {metas: true}
table_layout:
  default_order: [{field: indicador_id, order: DESC}]
```

## METAS
```yaml
table: metas
comment: Meta
tooltip: Meta de desempenho associada ao indicador e ao respetivo ano.
columns:
  meta_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da meta.", form_display: true, table_display: true, order: 1 }
  meta: { type: varchar, len: 255, nullable: false, comment: "Meta", tooltip: "Descrição da meta.", form_display: true, table_display: true, form_size: 8, order: 2 }
  meta_desc: { type: text, comment: "Descrição", tooltip: "Descrição detalhada da meta.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  ano_id: { type: integer, fk: "anos.ano_id", nullable: false, comment: "Ano", tooltip: "Ano da meta.", form_display: true, table_display: true, form_size: 4, order: 4 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano em formato literal.", form_display: true, table_display: true, form_size: 4, order: 5 }
  meta_valor: { type: decimal, len: 14, scale: 2, default: 0, comment: "Valor meta", tooltip: "Valor da meta para o ano.", form_display: true, table_display: true, form_size: 4, order: 6 }
  valor_atual: { type: decimal, len: 14, scale: 2, default: 0, comment: "Valor atual", tooltip: "Valor atual alcançado até ao momento.", form_display: true, table_display: true, form_size: 4, order: 7 }
  indicador_id: { type: integer, fk: "indicadores.indicador_id", nullable: false, comment: "Indicador", tooltip: "Indicador associado à meta.", form_display: true, table_display: true, form_size: 6, order: 8 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: meta_id, order: DESC}]
```

## ENTIDADES
```yaml
table: entidades
comment: Entidade
tooltip: Entidade pública, privada ou outra organização envolvida no programa.
columns:
  entidade_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da entidade.", form_display: true, table_display: true, order: 1 }
  entidade: { type: varchar, len: 255, nullable: false, comment: "Entidade", tooltip: "Nome da entidade.", form_display: true, table_display: true, form_size: 8, order: 2 }
  tipo_entidade_id: { type: integer, fk: "tipo_entidade.tipo_entidade_id", comment: "Tipo entidade", tooltip: "Tipo de entidade.", form_display: true, table_display: true, form_size: 4, order: 3 }
  entidade_pai_id: { type: integer, fk: "entidades.entidade_id", comment: "Entidade pai", tooltip: "Entidade de nível superior.", form_display: true, table_display: true, form_size: 4, order: 4 }
  concelho_id: { type: integer, fk: "concelhos.concelho_id", comment: "Concelho", tooltip: "Concelho da entidade.", form_display: true, table_display: true, form_size: 4, order: 5 }
  ilha_id: { type: integer, fk: "ilhas.ilha_id", comment: "Ilha", tooltip: "Ilha da entidade.", form_display: true, table_display: true, form_size: 4, order: 6 }
  genero_id: { type: integer, fk: "generos.genero_id", comment: "Género", tooltip: "Género ou categoria relevante da entidade.", form_display: true, table_display: true, form_size: 4, order: 7 }
  data_nasc_const: { type: date, comment: "Data nascimento/constituição", tooltip: "Data de nascimento ou constituição.", form_display: true, table_display: true, form_size: 4, order: 8 }
  email: { type: varchar, len: 255, comment: "Email", tooltip: "Contacto eletrónico.", form_display: true, table_display: true, form_size: 6, order: 9 }
  telefone: { type: varchar, len: 50, comment: "Telefone", tooltip: "Contacto telefónico.", form_display: true, table_display: true, form_size: 6, order: 10 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: entidade_id, order: DESC}]
```

## ANOS
```yaml
table: anos
comment: Ano
tooltip: Calendário de anos e respetivos períodos ativos para o programa.
columns:
  ano_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do ano.", form_display: true, table_display: true, order: 1 }
  ano: { type: varchar, len: 20, nullable: false, unique: true, comment: "Ano", tooltip: "Ano civil ou de referência.", form_display: true, table_display: true, form_size: 4, order: 2 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início do ano.", form_display: true, table_display: true, form_size: 4, order: 3 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim do ano.", form_display: true, table_display: true, form_size: 4, order: 4 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o ano está ativo.", form_display: true, table_display: true, form_size: 3, order: 5 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: ano, order: DESC}]
data:
  - {ano_id: 1, ano: "2025", data_ini: "2025-01-01", data_fim: "2025-12-31", activo: true, excluded: false}
  - {ano_id: 2, ano: "2026", data_ini: "2026-01-01", data_fim: "2026-12-31", activo: true, excluded: false}
  - {ano_id: 3, ano: "2027", data_ini: "2027-01-01", data_fim: "2027-12-31", activo: true, excluded: false}
```

## ILHAS
```yaml
table: ilhas
comment: Ilha
tooltip: Ilha geográfica de referência para entidades e regionalização.
columns:
  ilha_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da ilha.", form_display: true, table_display: true, order: 1 }
  ilha: { type: varchar, len: 100, nullable: false, unique: true, comment: "Ilha", tooltip: "Nome da ilha.", form_display: true, table_display: true, form_size: 6, order: 2 }
  ilha_desc: { type: text, comment: "Descrição", tooltip: "Descrição ou observação sobre a ilha.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se a ilha está ativa.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: ilha_id, order: ASC}]
data:
  - {ilha_id: 1, ilha: "Santo Antão", ilha_desc: "Ilha do norte do arquipélago", activo: true, excluded: false}
  - {ilha_id: 2, ilha: "São Vicente", ilha_desc: "Ilha do centro-norte", activo: true, excluded: false}
  - {ilha_id: 3, ilha: "Santa Luzia", ilha_desc: "Ilha pequena e não habitada", activo: true, excluded: false}
  - {ilha_id: 4, ilha: "São Nicolau", ilha_desc: "Ilha de orientação norte", activo: true, excluded: false}
  - {ilha_id: 5, ilha: "Sal", ilha_desc: "Ilha do arquipélago de Cabo Verde", activo: true, excluded: false}
  - {ilha_id: 6, ilha: "Boa Vista", ilha_desc: "Ilha com atividade turística", activo: true, excluded: false}
  - {ilha_id: 7, ilha: "Maio", ilha_desc: "Ilha do centro", activo: true, excluded: false}
  - {ilha_id: 8, ilha: "Santiago", ilha_desc: "Ilha principal do país", activo: true, excluded: false}
  - {ilha_id: 9, ilha: "Fogo", ilha_desc: "Ilha vulcânica", activo: true, excluded: false}
  - {ilha_id: 10, ilha: "Brava", ilha_desc: "Ilha mais pequena e montanhosa", activo: true, excluded: false}
```

## CONCELHOS
```yaml
table: concelhos
comment: Concelho
tooltip: Concelho geográfico para localização e regionalização das entidades.
columns:
  concelho_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do concelho.", form_display: true, table_display: true, order: 1 }
  ilha_id: { type: integer, fk: "ilhas.ilha_id", comment: "Ilha", tooltip: "Ilha a que o concelho pertence.", form_display: true, table_display: true, order: 2 }
  concelho: { type: varchar, len: 150, nullable: false, unique: true, comment: "Concelho", tooltip: "Nome do concelho.", form_display: true, table_display: true, form_size: 6, order: 3 }
  concelho_desc: { type: text, comment: "Descrição", tooltip: "Observações ou detalhe do concelho.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o concelho está ativo.", form_display: true, table_display: true, form_size: 3, order: 5 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: concelho_id, order: ASC}]
```

## GENEROS
```yaml
table: generos
comment: Género
tooltip: Classificação de género ou categoria de pessoa ou entidade.
columns:
  genero_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do género.", form_display: true, table_display: true, order: 1 }
  genero: { type: varchar, len: 100, nullable: false, unique: true, comment: "Género", tooltip: "Nome do género.", form_display: true, table_display: true, form_size: 6, order: 2 }
  genero_desc: { type: text, comment: "Descrição", tooltip: "Descrição ou observação do género.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o género está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: genero_id, order: ASC}]
data:
  - {genero_id: 1, genero: "Masculino", genero_desc: "Pessoa do sexo masculino", activo: true, excluded: false}
  - {genero_id: 2, genero: "Feminino", genero_desc: "Pessoa do sexo feminino", activo: true, excluded: false}
  - {genero_id: 3, genero: "Outro", genero_desc: "Outra categoria de identificação", activo: true, excluded: false}
```

## TIPO_CONTRATO
```yaml
table: tipo_contrato
comment: Tipo de Contrato
tooltip: Tipologia de contratos e modalidades de acordo com o programa.
columns:
  tipo_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do tipo de contrato.", form_display: true, table_display: true, order: 1 }
  tipo_contrato: { type: varchar, len: 150, nullable: false, unique: true, comment: "Tipo", tooltip: "Nome do tipo de contrato.", form_display: true, table_display: true, form_size: 6, order: 2 }
  tipo_contrato_desc: { type: text, comment: "Descrição", tooltip: "Descrição da tipologia do contrato.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o tipo está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: tipo_contrato_id, order: ASC}]
```

## TIPO_ENTIDADE
```yaml
table: tipo_entidade
comment: Tipo de Entidade
tooltip: Tipo ou categoria da entidade envolvida na implementação do programa.
columns:
  tipo_entidade_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do tipo de entidade.", form_display: true, table_display: true, order: 1 }
  tipo_entidade: { type: varchar, len: 150, nullable: false, unique: true, comment: "Tipo", tooltip: "Nome do tipo de entidade.", form_display: true, table_display: true, form_size: 6, order: 2 }
  tipo_entidade_desc: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de entidade.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o tipo está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: tipo_entidade_id, order: ASC}]
```

## STATUS_ACAO
```yaml
table: status_acao
comment: Status da Ação
tooltip: Estados possíveis para acompanhamento do progresso da ação.
columns:
  status_acao_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do status da ação.", form_display: true, table_display: true, order: 1 }
  status_acao: { type: varchar, len: 100, nullable: false, unique: true, comment: "Status", tooltip: "Nome do estado.", form_display: true, table_display: true, form_size: 6, order: 2 }
  status_acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o estado está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: status_acao_id, order: ASC}]
data:
  - {status_acao_id: 1, status_acao: "Planeada", status_acao_desc: "Ação ainda em planeamento", activo: true, excluded: false}
  - {status_acao_id: 2, status_acao: "Em execução", status_acao_desc: "Ação em curso", activo: true, excluded: false}
  - {status_acao_id: 3, status_acao: "Concluída", status_acao_desc: "Ação concluída", activo: true, excluded: false}
  - {status_acao_id: 4, status_acao: "Suspensa", status_acao_desc: "Ação suspensa", activo: true, excluded: false}
```

## STATUS_CONTRATO
```yaml
table: status_contrato
comment: Status do Contrato
tooltip: Estados possíveis para acompanhamento de contratos.
columns:
  status_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do status do contrato.", form_display: true, table_display: true, order: 1 }
  status_contrato: { type: varchar, len: 100, nullable: false, unique: true, comment: "Status", tooltip: "Nome do estado do contrato.", form_display: true, table_display: true, form_size: 6, order: 2 }
  status_contrato_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado do contrato.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o status está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: status_contrato_id, order: ASC}]
```

## STATUS_EXECUCAO
```yaml
table: status_execucao
comment: Status de Execução
tooltip: Estados possíveis para fases de execução do contrato.
columns:
  status_execucao_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do status de execução.", form_display: true, table_display: true, order: 1 }
  status_execucao: { type: varchar, len: 100, nullable: false, unique: true, comment: "Status", tooltip: "Nome do estado de execução.", form_display: true, table_display: true, form_size: 6, order: 2 }
  status_execucao_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado da execução.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o estado está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: status_execucao_id, order: ASC}]
```

## UNIDADES
```yaml
table: unidades
comment: Unidades
tooltip: Unidades de medida e referência para indicadores e cálculos.
columns:
  unidade_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da unidade.", form_display: true, table_display: true, order: 1 }
  unidade: { type: varchar, len: 100, nullable: false, unique: true, comment: "Unidade", tooltip: "Nome da unidade de medida.", form_display: true, table_display: true, form_size: 6, order: 2 }
  unidade_desc: { type: text, comment: "Descrição", tooltip: "Descrição da unidade.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  simbolo: { type: varchar, len: 20, comment: "Símbolo", tooltip: "Símbolo ou abreviatura da unidade.", form_display: true, table_display: true, form_size: 3, order: 4 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se a unidade está ativa.", form_display: true, table_display: true, form_size: 3, order: 5 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: unidade_id, order: ASC}]
data:
  - {unidade_id: 1, unidade: "Percentagem", unidade_desc: "Percentagem relativa", simbolo: "%", activo: true, excluded: false}
  - {unidade_id: 2, unidade: "Número", unidade_desc: "Quantidade absoluta", simbolo: "#", activo: true, excluded: false}
  - {unidade_id: 3, unidade: "Montante", unidade_desc: "Valor monetário", simbolo: "MVA", activo: true, excluded: false}
  - {unidade_id: 4, unidade: "Dias", unidade_desc: "Duração em dias", simbolo: "dias", activo: true, excluded: false}
```
