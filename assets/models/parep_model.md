<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# PAREP_MODEL
```yaml
name: PAREP
description: Modelo de gestão do programa PAREP-CV para objetivos, resultados esperados, indicadores, ações, contratos, entidades e parametrização.
runs_as: MODEL
database: PAREP
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    active: true
    menu_config: '{"label": "dashboard","tooltip": "dashboard_desc","load_items": {"table": "dashboard","tables": ["dashboard"]}}'
    tables:
      - dashboard
  PAREP:
    menu_icon: chart-pie
    menu_order: 1
    active: true
    menu_config: |
      {
        "label": "objetivos",
        "tooltip": "desc_objetivo",
        "load_items": {
          "table": "objetivos",
          "label": "objetivo",
          "tooltip": "desc_objetivo",
          "detail": true,
          "load_items": [
            {
              "table": "resultados",
              "label": "resultado",
              "tooltip": "desc_resultado",
              "detail": true,
              "load_items": {
                "table": "resultado_anos",
                "label": "ano",
                "tooltip": "ano",
                "detail": true,
                "load_items": {
                  "table": "acoes",
                  "label": "acao",
                  "tooltip": "acao_desc",
                  "detail": false
                }
              }
            },
            {
              "table": "indicadores",
              "label": "indicador",
              "tooltip": "desc_indicador",
              "detail": false
            }
          ],
          "tables": ["objetivos", "resultados", "indicadores", "resultado_anos", "acoes"]
        }
      }
    tables:
      - objetivos
      - {table: resultados, active: false}
      - {table: resultado_anos, active: false}
      - {table: acoes, active: false}
      - {table: contratos, active: false}
      - {table: contrato_execucoes, active: false}
      - {table: contrato_anexos, active: false}
      - {table: logistica, active: false}
      - {table: beneficiarios, active: false}
      - {table: indicadores, active: false}
      - {table: metas, active: false}
  Parametrização:
    menu_icon: adjustments
    menu_order: 2
    active: true
    tables:
      - entidades
      - ambitos
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

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar, len: 200, comment: "Dashboard", form_display: true, table_display: true, form_size: 8, order: 1 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, order: 5 }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 2, order: 2 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 3 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - { component: EvidenceDash, label: dashboard, intercept_r: true, size: 12 }
```

## AMBITOS
```yaml
table: ambitos
comment: Âmbito
tooltip: Classificação de âmbito para categorizar objetivos e resultados (Pré-escolar, 1º Ciclo, etc).
columns:
  ambito_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do âmbito.", form_display: true, table_display: true, order: 1 }
  ambito: { type: varchar, len: 150, nullable: false, unique: true, comment: "Âmbito", tooltip: "Nome do âmbito.", form_display: true, table_display: true, form_size: 6, order: 2 }
  ambito_desc: { type: text, comment: "Descrição", tooltip: "Descrição do âmbito.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  activo: { type: boolean, default: true, comment: "Ativo", tooltip: "Indica se o âmbito está ativo.", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: ambito_id, order: ASC}]
data:
  - {ambito_id: 1, ambito: "Pré-escolar", ambito_desc: "Educação Pré-escolar de qualidade, acessível e inclusiva", activo: true, excluded: false}
  - {ambito_id: 2, ambito: "1º Ciclo", ambito_desc: "Ensino Básico de qualidade, equitativo, inclusivo e com êxito educativo", activo: true, excluded: false}
```

## OBJETIVOS
```yaml
table: objetivos
comment: Objetivo
tooltip: Objetivo estratégico do programa PAREP-CV.
columns:
  objetivo_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador único do objetivo.", form_display: true, table_display: true, order: 1 }
  objetivo: { type: varchar, len: 255, nullable: false, comment: "Objetivo", tooltip: "Nome do objetivo estratégico.", form_display: true, table_display: true, form_size: 8, order: 2 }
  desc_objetivo: { type: text, comment: "Descrição", tooltip: "Descrição detalhada do objetivo.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  ambito_id: { type: integer, fk: "ambitos.ambito_id", nullable: false, comment: "Âmbito", tooltip: "Âmbito do objetivo (Pré-escolar, 1º Ciclo, etc).", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início do período de execução do objetivo.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de término do período de execução do objetivo.", form_display: true, table_display: true, form_size: 4, order: 6 }
  status_id: { type: integer, comment: "Estado", tooltip: "Estado do objetivo.", form_display: true, table_display: true, form_size: 4, order: 7 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {resultados: true, indicadores: true}
table_layout:
  default_order: [{field: objetivo_id, order: DESC}]
```

## RESULTADOS
```yaml
table: resultados
comment: Resultado Esperado
tooltip: Resultado esperado associado a um objetivo.
columns:
  resultado_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do resultado esperado.", form_display: true, table_display: true, order: 1 }
  resultado: { type: varchar, len: 255, nullable: false, comment: "Resultado Esperado", tooltip: "Nome do resultado esperado.", form_display: true, table_display: true, form_size: 8, order: 2 }
  desc_resultado: { type: text, comment: "Descrição", tooltip: "Descrição do resultado esperado.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo ao qual o resultado esperado está relacionado.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da execução.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da execução.", form_display: true, table_display: true, form_size: 4, order: 6 }
  status_id: { type: integer, comment: "Estado", tooltip: "Estado do resultado esperado.", form_display: true, table_display: true, form_size: 4, order: 7 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {resultado_anos: true}
table_layout:
  default_order: [{field: resultado_id, order: DESC}]
```

## RESULTADO_ANOS
```yaml
table: resultado_anos
comment: resultado Ano
tooltip: Variação anual de um resultado esperado associado ao ano e respetiva execução.
columns:
  resultado_ano_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do registo anual do resultado esperado.", form_display: true, table_display: true, order: 1 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", nullable: false, comment: "Resultado Esperado", tooltip: "Resultado esperado principal.", form_display: true, table_display: true, form_size: 4, order: 2 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo associado.", form_display: true, table_display: true, form_size: 4, order: 3 }
  ano_id: { type: integer, fk: "anos.ano_id", nullable: false, comment: "Ano", tooltip: "Ano de execução do resultado esperado.", form_display: true, table_display: true, form_size: 4, order: 4 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano literal para referência rápida.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da execução anual.", form_display: true, table_display: true, form_size: 4, order: 6 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da execução anual.", form_display: true, table_display: true, form_size: 4, order: 7 }
  meta_ano: { type: decimal, len: 14, scale: 2, default: 0, comment: "Meta anual", tooltip: "Meta de execução anual do resultado esperado.", form_display: true, table_display: true, form_size: 4, order: 8 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  default_order: [{field: resultado_ano_id, order: DESC}]
```

## ACOES
```yaml
table: acoes
comment: Ação
tooltip: Ação de implementação associada a um resultado esperado anual.
columns:
  acao_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da ação.", form_display: true, table_display: true, order: 1 }
  acao: { type: varchar, len: 255, nullable: false, comment: "Ação", tooltip: "Nome da ação.", form_display: true, table_display: true, form_size: 8, order: 2 }
  acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição da ação.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  data_ini: { type: date, comment: "Data início", tooltip: "Data de início da ação.", form_display: true, table_display: true, form_size: 4, order: 4 }
  data_fim: { type: date, comment: "Data fim", tooltip: "Data de fim da ação.", form_display: true, table_display: true, form_size: 4, order: 5 }
  orcamento_estimado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento estimado", tooltip: "Valor estimado da ação.", form_display: true, table_display: true, form_size: 4, order: 6 }
  orcamento_executado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento executado", tooltip: "Valor executado da ação.", form_display: true, table_display: true, form_size: 4, order: 7 }
  status_id: { type: integer, fk: "status_acao.status_acao_id", comment: "Status", tooltip: "Estado atual da ação.", form_display: true, table_display: true, form_size: 4, order: 8 }
  resultado_ano_id: { type: integer, fk: "resultado_anos.resultado_ano_id", comment: "Resultado Esperado Ano", tooltip: "Relação com a execução anual do resultado esperado.", form_display: true, table_display: true, form_size: 4, order: 9 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano da ação for referência.", form_display: true, table_display: true, form_size: 4, order: 10 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado Esperado", tooltip: "Resultado esperado associado.", form_display: true, table_display: true, form_size: 4, order: 11 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo", tooltip: "Objetivo associado.", form_display: true, table_display: true, form_size: 4, order: 12 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  origem_id: { type: integer, fk: "contratos.contrato_id", comment: "Contrato origem / pai", tooltip: "Contrato principal do qual este contrato deriva ou é adenda.", form_display: true, table_display: true, form_size: 4, order: 9 }
  addenda_id: { type: integer, fk: "contratos.contrato_id", comment: "Adenda / contrato destino", tooltip: "Contrato de adenda ou contrato destino relacionado com este contrato principal.", form_display: true, table_display: true, form_size: 4, order: 10 }
  resultado_ano_id: { type: integer, fk: "resultado_anos.resultado_ano_id", comment: "Resultado Esperado Ano", tooltip: "Relação com o resultado esperado anual.", form_display: true, table_display: true, form_size: 4, order: 11 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano do contrato.", form_display: true, table_display: true, form_size: 4, order: 12 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado Esperado", tooltip: "Resultado esperado associado.", form_display: true, table_display: true, form_size: 4, order: 13 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo", tooltip: "Objetivo associado.", form_display: true, table_display: true, form_size: 4, order: 14 }
  status_id: { type: integer, fk: "status_contrato.status_contrato_id", comment: "Status", tooltip: "Estado do contrato.", form_display: true, table_display: true, form_size: 4, order: 15 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
tooltip: Gestão logística associada à resultado anual e respetivo acompanhamento de custos.
columns:
  logistica_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador da logística.", form_display: true, table_display: true, order: 1 }
  logistica: { type: varchar, len: 255, nullable: false, comment: "Logística", tooltip: "Descrição da logística ou operação.", form_display: true, table_display: true, form_size: 8, order: 2 }
  orcamento_estimado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento estimado", tooltip: "Estimativa de custos da logística.", form_display: true, table_display: true, form_size: 4, order: 3 }
  orcamento_executado: { type: decimal, len: 14, scale: 2, default: 0, comment: "Orçamento executado", tooltip: "Valor executado da logística.", form_display: true, table_display: true, form_size: 4, order: 4 }
  status_id: { type: integer, fk: "status_acao.status_acao_id", comment: "Status", tooltip: "Estado da logística.", form_display: true, table_display: true, form_size: 4, order: 5 }
  attach_logistica_anexo: { type: varchar, len: 255, comment: "Anexo", tooltip: "Anexo ou documento relacionado com a logística.", form_display: true, table_display: true, form_size: 6, order: 6 }
  acao_id: { type: integer, fk: "acaos.acao_id", comment: "Resultado Esperado Ano", tooltip: "Resultado esperado anual associado.", form_display: true, table_display: true, form_size: 6, order: 7 }
  resultado_ano_id: { type: integer, fk: "resultado_anos.resultado_ano_id", comment: "Resultado Esperado Ano", tooltip: "Resultado esperado anual associado.", form_display: true, table_display: true, form_size: 6, order: 7 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano da logística.", form_display: true, table_display: true, form_size: 4, order: 8 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado Esperado", tooltip: "Resultado esperado principal.", form_display: true, table_display: true, form_size: 4, order: 9 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo", tooltip: "Objetivo associado.", form_display: true, table_display: true, form_size: 4, order: 10 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
tooltip: Registo dos beneficiários associados a uma resultado anual.
columns:
  beneficiario_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Identificador do beneficiário.", form_display: true, table_display: true, order: 1 }
  beneficiario: { type: varchar, len: 255, nullable: false, comment: "Beneficiário", tooltip: "Nome do beneficiário.", form_display: true, table_display: true, form_size: 8, order: 2 }
  data_ini_benef: { type: date, comment: "Data início", tooltip: "Data de início da elegibilidade ou apoio.", form_display: true, table_display: true, form_size: 4, order: 3 }
  data_fim_benef: { type: date, comment: "Data fim", tooltip: "Data de fim do apoio ou acompanhamento.", form_display: true, table_display: true, form_size: 4, order: 4 }
  beneficiario_status_id: { type: integer, comment: "Status", tooltip: "Status do beneficiário.", form_display: true, table_display: true, form_size: 4, order: 5 }
  attach_beneficiarios: { type: varchar, len: 255, comment: "Anexo", tooltip: "Anexo ou ficheiro associado ao beneficiário.", form_display: true, table_display: true, form_size: 6, order: 6 }
  resultado_ano_id: { type: integer, fk: "resultado_anos.resultado_ano_id", comment: "Resultado Esperado Ano", tooltip: "Resultado esperado anual associado.", form_display: true, table_display: true, form_size: 6, order: 7 }
  ano: { type: varchar, len: 20, comment: "Ano", tooltip: "Ano de referência.", form_display: true, table_display: true, form_size: 4, order: 8 }
  resultado_id: { type: integer, fk: "resultados.resultado_id", comment: "Resultado Esperado", tooltip: "Resultado esperado principal.", form_display: true, table_display: true, form_size: 4, order: 9 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo", tooltip: "Objetivo associado.", form_display: true, table_display: true, form_size: 4, order: 10 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo ao qual o indicador pertence.", form_display: true, table_display: true, form_size: 6, order: 9 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  desc_entidade: { type: text, comment: "Descrição", tooltip: "Descrição geral da entidade, missão, função ou observações relevantes.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 3 }
  tipo_entidade_id: { type: integer, fk: "tipo_entidade.tipo_entidade_id", comment: "Tipo entidade", tooltip: "Tipo de entidade.", form_display: true, table_display: true, form_size: 4, order: 4 }
  entidade_pai_id: { type: integer, fk: "entidades.entidade_id", comment: "Entidade pai", tooltip: "Entidade de nível superior.", form_display: true, table_display: true, form_size: 4, order: 5 }
  concelho_id: { type: integer, fk: "concelhos.concelho_id", comment: "Concelho", tooltip: "Concelho da entidade.", form_display: true, table_display: true, form_size: 4, order: 6 }
  ilha_id: { type: integer, fk: "ilhas.ilha_id", comment: "Ilha", tooltip: "Ilha da entidade.", form_display: true, table_display: true, form_size: 4, order: 7 }
  genero_id: { type: integer, fk: "generos.genero_id", comment: "Género", tooltip: "Género ou categoria relevante da entidade.", form_display: true, table_display: true, form_size: 4, order: 8 }
  data_nasc_const: { type: date, comment: "Data nascimento/constituição", tooltip: "Data de nascimento ou constituição.", form_display: true, table_display: true, form_size: 4, order: 9 }
  email: { type: varchar, len: 255, comment: "Email", tooltip: "Contacto eletrónico.", form_display: true, table_display: true, form_size: 6, order: 10 }
  telefone: { type: varchar, len: 50, comment: "Telefone", tooltip: "Contacto telefónico.", form_display: true, table_display: true, form_size: 6, order: 11 }
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  - {ilha_id: 6, ilha: "Boa Vista", ilha_desc: "Ilha com resultado turística", activo: true, excluded: false}
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: concelho_id, order: ASC}]
data:
  - {concelho_id: 1, ilha_id: 1, concelho: "Ribeira Grande", concelho_desc: "Concelho de Santo Antão", activo: true, excluded: false}
  - {concelho_id: 2, ilha_id: 1, concelho: "Paul", concelho_desc: "Concelho de Santo Antão", activo: true, excluded: false}
  - {concelho_id: 3, ilha_id: 1, concelho: "Porto Novo", concelho_desc: "Concelho de Santo Antão", activo: true, excluded: false}
  - {concelho_id: 4, ilha_id: 2, concelho: "São Vicente", concelho_desc: "Concelho de São Vicente", activo: true, excluded: false}
  - {concelho_id: 5, ilha_id: 4, concelho: "Ribeira Brava", concelho_desc: "Concelho de São Nicolau", activo: true, excluded: false}
  - {concelho_id: 6, ilha_id: 4, concelho: "Tarrafal de São Nicolau", concelho_desc: "Concelho de São Nicolau", activo: true, excluded: false}
  - {concelho_id: 7, ilha_id: 5, concelho: "Sal", concelho_desc: "Concelho do Sal", activo: true, excluded: false}
  - {concelho_id: 8, ilha_id: 6, concelho: "Boa Vista", concelho_desc: "Concelho da Boa Vista", activo: true, excluded: false}
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: tipo_entidade_id, order: ASC}]
data:
  - {tipo_entidade_id: 1, tipo_entidade: "Tipo", tipo_entidade_desc: "Entidade genérica ou de natureza institucional", activo: true, excluded: false}
  - {tipo_entidade_id: 2, tipo_entidade: "Consultor", tipo_entidade_desc: "Entidade ou profissional prestador de serviços de consultoria", activo: true, excluded: false}
  - {tipo_entidade_id: 3, tipo_entidade: "Empresa fornecedora", tipo_entidade_desc: "Empresa que fornece bens ou serviços ao programa", activo: true, excluded: false}
  - {tipo_entidade_id: 4, tipo_entidade: "Empresa construção", tipo_entidade_desc: "Empresa especializada em obras, construção ou instalação", activo: true, excluded: false}
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: status_contrato_id, order: ASC}]
data:
  - {status_contrato_id: 1, status_contrato: "Rascunho", status_contrato_desc: "Contrato em elaboração", activo: true, excluded: false}
  - {status_contrato_id: 2, status_contrato: "Em negociação", status_contrato_desc: "Contrato em processo de negociação", activo: true, excluded: false}
  - {status_contrato_id: 3, status_contrato: "Assinado", status_contrato_desc: "Contrato formalizado e assinado", activo: true, excluded: false}
  - {status_contrato_id: 4, status_contrato: "Em execução", status_contrato_desc: "Contrato em execução", activo: true, excluded: false}
  - {status_contrato_id: 5, status_contrato: "Concluído", status_contrato_desc: "Contrato concluído com êxito", activo: true, excluded: false}
  - {status_contrato_id: 6, status_contrato: "Cancelado", status_contrato_desc: "Contrato cancelado ou não executado", activo: true, excluded: false}
  - {status_contrato_id: 7, status_contrato: "Suspenso", status_contrato_desc: "Contrato suspenso temporariamente", activo: true, excluded: false}
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização.", form_display: false, table_display: true }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o registo está excluído.", form_display: false, table_display: false }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: status_execucao_id, order: ASC}]
data:
  - {status_execucao_id: 1, status_execucao: "Draft", status_execucao_desc: "Registo em rascunho", activo: true, excluded: false}
  - {status_execucao_id: 2, status_execucao: "Iniciado", status_execucao_desc: "Execução iniciada", activo: true, excluded: false}
  - {status_execucao_id: 3, status_execucao: "Em execução", status_execucao_desc: "Execução em progresso", activo: true, excluded: false}
  - {status_execucao_id: 4, status_execucao: "Cancelado", status_execucao_desc: "Execução cancelada", activo: true, excluded: false}
  - {status_execucao_id: 5, status_execucao: "Adendado", status_execucao_desc: "Execução alterada por adenda", activo: true, excluded: false}
  - {status_execucao_id: 6, status_execucao: "Concluído", status_execucao_desc: "Execução concluída", activo: true, excluded: false}
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
  user_id: { type: integer, comment: "Utilizador", tooltip: "Utilizador responsável pelo registo.", form_display: false, table_display: false }
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

<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# PAREP_DATA
```yaml
name: PAREP_DATA
description: Dados UI do PAREP Compact - Objetivos, Resultados Esperados, Indicadores e Metas 2025-2026 por Âmbito
database: PAREP
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## DASHBOARD1
```yaml
table: dashboard
description: Add default Dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id:   1
  dashboard:      Resumo Indicadores
  dashboard_desc: Exemplo Layout Resumo Indicadores
  dashboard_conf: FileContent(parep/dashboard1.md)
  order:          1
  active:         true
```

## OBJETIVO1
```yaml
table: objetivos
description: Add OBJ 1
cond: "where objetivo_id = :objetivo_id"
data:
  objetivo_id: 1
  objetivo: OBJ1
  desc_objetivo: Assegurar que todas as crianças de 4 e 5 anos frequentem o Pré-escolar
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O1R1
          desc_resultado: Acesso universal e inclusivo à Educação Pré-escolar.
          objetivo_id: objetivo_id()
        - resultado: O1R2
          desc_resultado: Redução das assimetrias regionais de acesso.
          objetivo_id: objetivo_id()
        - resultado: O1R3
          desc_resultado: Aumento da frequência das crianças de 4 e 5 anos nos jardins de infância.
          objetivo_id: objetivo_id()
        - resultado: O1R4
          desc_resultado: Maior inclusão de crianças provenientes de famílias vulneráveis.
          objetivo_id: objetivo_id()    
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O1IND1
          desc_indicador: Taxa líquida de escolarização/cobertura das crianças de 4-5 anos no Pré-escolar
          valor_baseline: 0.841
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
                meta: O1IND12026
                ano_id: 1
                ano: 2026
                meta_valor: 0.95
                indicador_id: indicador_id()
        - indicador: O1IND2
          desc_indicador: Nº de crianças que frequentam Jardins de Infância
          valor_baseline: 15906
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
                meta: O1IND22026
                ano_id: 1
                ano: 2026
                meta_valor: 17280
                indicador_id: indicador_id()
        - indicador: O1IND3
          desc_indicador: Nº de famílias de classe 1 e 2 que recebem apoio para escolarização dos seus filhos na Educação Pré-escolar
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
                meta: O1IND32026
                ano_id: 1
                ano: 2026
                meta_valor: 3500
                indicador_id: indicador_id()
        - indicador: O1IND4
          desc_indicador: Nº de Jardins que recebem kits de materiais lúdico-pedagógicos e equipamentos para crianças com NEE
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
                meta: O1IND342026
                ano_id: 1
                ano: 2026
                meta_valor: 150
                indicador_id: indicador_id()
```

## OBJETIVO2
```yaml
table: objetivos
cond: "where objetivo_id = :objetivo_id"
data:
  objetivo_id: 2
  objetivo: OBJ2
  desc_objetivo: Melhorar o desempenho das aprendizagens das crianças no Pré-escolar
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O2R1
          desc_resultado: Melhoria das competências em leitura, escrita e numeracia.
          objetivo_id: objetivo_id()
        - resultado: O2R2
          desc_resultado: Profissionais de infância mais qualificados.
          objetivo_id: objetivo_id()
        - resultado: O2R3
          desc_resultado: Contextos pedagógicos de aprendizagem melhorados.
          objetivo_id: objetivo_id()
        - resultado: O2R4
          desc_resultado: Práticas pedagógicas alinhadas com as orientações curriculares.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O2IND5
          desc_indicador: Percentagem de crianças de 4-5 anos que desenvolvem competências básicas em língua portuguesa e conhecimento dos números
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
        - indicador: O2IND6
          desc_indicador: Percentagem de JI que utilizam e cumprem o programa previsto em articulação com as unidades educativas
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
        - indicador: O2IND7
          desc_indicador: Número de educadoras com formação inicial formadas e integradas no sistema
          valor_baseline: 136
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O2IND72026
              ano_id: 1
              ano: 2026
              meta_valor: 424
              indicador_id: indicador_id()
        - indicador: O2IND8
          desc_indicador: Número de monitoras com formação inicial formadas e integradas no sistema
          valor_baseline: 283
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O2IND82026
              ano_id: 1
              ano: 2026
              meta_valor: 566
              indicador_id: indicador_id()
        - indicador: O2IND9A
          desc_indicador: Percentagem de educadoras com formação adequada
          valor_baseline: 0.096
          unidades_id: 1
          objetivo_id: objetivo_id()
        - indicador: O2IND9B
          desc_indicador: Percentagem de monitoras com formação adequada
          valor_baseline: 0.20
          unidades_id: 1
          objetivo_id: objetivo_id()
```

## OBJETIVO3
```yaml
table: objetivos
cond: "where objetivo_id = :objetivo_id"
data:
  objetivo_id: 3
  objetivo: OBJ3
  desc_objetivo: Melhorar a eficiência e eficácia do uso dos recursos disponibilizados ao Pré-escolar
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O3R1
          desc_resultado: Planeamento mais eficiente da rede pré-escolar.
          objetivo_id: objetivo_id()
        - resultado: O3R2
          desc_resultado: Melhor utilização dos recursos humanos e infraestruturas.
          objetivo_id: objetivo_id()
        - resultado: O3R3
          desc_resultado: Melhor gestão dos recursos disponibilizados ao subsistema.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O3IND10
          desc_indicador: Carta Escolar do Pré-escolar elaborada
          valor_baseline: null
          unidades_id: 3
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O3IND102026
              ano_id: 1
              ano: 2026
              meta_valor: 1
              indicador_id: indicador_id()
        - indicador: O3IND11
          desc_indicador: Rácio criança/profissional de infância
          valor_baseline: 20
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O3IND112026
              ano_id: 1
              ano: 2026
              meta_valor: 25
              indicador_id: indicador_id()
```

## OBJETIVO4
```yaml
table: objetivos
cond: "where objetivo_id = :objetivo_id"
data:
  objetivo_id: 4
  objetivo: OBJ4
  desc_objetivo: Reforçar a capacidade institucional e organizativa da Educação Pré-escolar
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O4R1
          desc_resultado: Novo regime jurídico da Educação Pré-escolar implementado.
          objetivo_id: objetivo_id()
        - resultado: O4R2
          desc_resultado: Estatuto de carreira dos profissionais de infância definido e operacionalizado.
          objetivo_id: objetivo_id()
        - resultado: O4R3
          desc_resultado: Gestão dos jardins de infância baseada em resultados.
          objetivo_id: objetivo_id()
        - resultado: O4R4
          desc_resultado: Maior articulação entre jardins de infância e escolas básicas.
          objetivo_id: objetivo_id()
        - resultado: O4R5
          desc_resultado: Utilização do SIGE para monitoria e gestão do subsistema.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O4IND12
          desc_indicador: Percentagem de JI que cumprem os requisitos mínimos exigidos para funcionamento
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND122026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
        - indicador: O4IND13
          desc_indicador: Nº de Gestores de Jardins de Infância capacitados em gestão baseada em resultados
          valor_baseline: 582
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND132026
              ano_id: 1
              ano: 2026
              meta_valor: 600
              indicador_id: indicador_id()
        - indicador: O4IND14
          desc_indicador: Percentagem de gestores dos JI habilitados para a função de direção
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND142026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
        - indicador: O4IND15
          desc_indicador: Nº de JI que fazem intercâmbio e articulam com as escolas básicas da sua área de influência
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND152026
              ano_id: 1
              ano: 2026
              meta_valor: 360
              indicador_id: indicador_id()
        - indicador: O4IND16
          desc_indicador: Estatuto de carreira dos profissionais de infância elaborado discutido e apropriado
          valor_baseline: null
          unidades_id: 3
          objetivo_id: objetivo_id()
        - indicador: O4IND17
          desc_indicador: Nº de Jardins com prática de gestão baseada em resultados
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND172026
              ano_id: 1
              ano: 2026
              meta_valor: 600
              indicador_id: indicador_id()
        - indicador: O4IND18
          desc_indicador: Percentagem de Jardins com prática de gestão baseada em resultados
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND182026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
        - indicador: O4IND19
          desc_indicador: Nº de Gestores de Jardins capacitados para uso e manuseamento do SIGE
          valor_baseline: 582
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND192026
              ano_id: 1
              ano: 2026
              meta_valor: 600
              indicador_id: indicador_id()
        - indicador: O4IND20
          desc_indicador: Percentagem de JI que disponibiliza os seus dados através do SIGE
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O4IND202026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
```

## OBJETIVO5
```yaml
table: objetivos
cond: where objetivo_id = :objetivo_id
data:
  objetivo_id: 5
  objetivo: OBJ5
  desc_objetivo: Consolidar o acesso equitativo e inclusivo no 1.º Ciclo do EBO
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O5R1
          desc_resultado: Acesso universal e inclusivo consolidado no EBO.
          objetivo_id: objetivo_id()
        - resultado: O5R2
          desc_resultado: Melhor integração das crianças com NEE.
          objetivo_id: objetivo_id()
        - resultado: O5R3
          desc_resultado: Redução das desigualdades de género e vulnerabilidade social.
          objetivo_id: objetivo_id()
        - resultado: O5R4
          desc_resultado: Reforço da capacidade de resposta das estruturas de educação inclusiva.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O5IND21
          desc_indicador: Nº de campanhas regulares de sensibilização nas comunidades implementadas visando reforçar o acesso e reduzir o risco de abandono escolar
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O5IND212026
              ano_id: 1
              ano: 2026
              meta_valor: 20
              indicador_id: indicador_id()
        - indicador: O5IND22
          desc_indicador: Estudo de mapeamento de crianças com NEE elaborado
          valor_baseline: null
          unidades_id: 3
          objetivo_id: objetivo_id()
        - indicador: O5IND23
          desc_indicador: Equipas EMAEI reforçadas com mais um técnico em todos os concelhos do país
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O5IND232026
              ano_id: 1
              ano: 2026
              meta_valor: 22
              indicador_id: indicador_id()
        - indicador: O5IND24
          desc_indicador: Reforço dos materiais lúdico-pedagógicos destinados às crianças com NEE
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O5IND242026
              ano_id: 1
              ano: 2026
              meta_valor: 22
              indicador_id: indicador_id()
        - indicador: O5IND25
          desc_indicador: Taxa líquida de escolarização das crianças de 6-13 anos no Ensino Básico Obrigatório
          valor_baseline: 0.996
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O5IND252026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
        - indicador: O5IND26
          desc_indicador: Índice de Paridade de Género (M/F)
          valor_baseline: 0.91
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O5IND262026
              ano_id: 1
              ano: 2026
              meta_valor: 0.95
              indicador_id: indicador_id()
```

## OBJETIVO6
```yaml
table: objetivos
cond: "where objetivo_id = :objetivo_id"
data:
  objetivo_id: 6
  objetivo: OBJ6
  desc_objetivo: Reforçar o êxito escolar e a qualidade das aprendizagens em Português e Matemática
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O6R1
          desc_resultado: Melhoria do desempenho dos alunos em Língua Portuguesa.
          objetivo_id: objetivo_id()
        - resultado: O6R2
          desc_resultado: Melhoria do desempenho dos alunos em Matemática.
          objetivo_id: objetivo_id()
        - resultado: O6R3
          desc_resultado: Professores mais capacitados na gestão curricular e avaliação.
          objetivo_id: objetivo_id()
        - resultado: O6R4
          desc_resultado: Maior inclusão das crianças com NEE no processo de ensino-aprendizagem.
          objetivo_id: objetivo_id()
        - resultado: O6R5
          desc_resultado: Aumento das taxas de sucesso escolar.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O6IND27
          desc_indicador: Nº de docentes do 1.º ciclo que recebem ações de capacitação em gestão curricular, técnicas de avaliação e diferenciação pedagógica
          valor_baseline: 4188
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND272026
              ano_id: 1
              ano: 2026
              meta_valor: 4200
              indicador_id: indicador_id()
        - indicador: O6IND28
          desc_indicador: Nº de docentes capacitados em matéria de gestão pedagógica de crianças com NEE
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND282026
              ano_id: 1
              ano: 2026
              meta_valor: 4200
              indicador_id: indicador_id()
        - indicador: O6IND29
          desc_indicador: Nº de equipas pedagógicas nacional e concelhias capacitadas para reforço da supervisão, acompanhamento e avaliação no EBO
          valor_baseline: 60
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND292026
              ano_id: 1
              ano: 2026
              meta_valor: 60
              indicador_id: indicador_id()
        - indicador: O6IND30
          desc_indicador: Percentagem de aprovação no 1.º Ciclo do EBO
          valor_baseline: 0.839
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND302026
              ano_id: 1
              ano: 2026
              meta_valor: 0.92
              indicador_id: indicador_id()
        - indicador: O6IND31
          desc_indicador: Percentagem de reprovação no 1.º Ciclo do EBO
          valor_baseline: 0.152
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND312026
              ano_id: 1
              ano: 2026
              meta_valor: 0.08
              indicador_id: indicador_id()
        - indicador: O6IND32
          desc_indicador: Percentagem de abandono no 1.º Ciclo do EBO
          valor_baseline: 0.0001
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND322026
              ano_id: 1
              ano: 2026
              meta_valor: 0.00
              indicador_id: indicador_id()
        - indicador: O6IND33
          desc_indicador: Percentagem de crianças do EBO que adquirem competências básicas em língua portuguesa (leitura, escrita) e matemática
          valor_baseline: 0.327
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND332026
              ano_id: 1
              ano: 2026
              meta_valor: 0.60
              indicador_id: indicador_id()
        - indicador: O6IND34
          desc_indicador: Percentagem de docentes abrangidos pelo programa de formação contínua de acordo com o plano de formação de professores
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O6IND342026
              ano_id: 1
              ano: 2026
              meta_valor: 0.90
              indicador_id: indicador_id()
```

## OBJETIVO7
```yaml
table: objetivos
cond: where objetivo_id = :objetivo_id
data:
  objetivo_id: 7
  objetivo: OBJ7
  desc_objetivo: Melhorar a eficiência e eficácia do uso dos recursos disponibilizados ao EBO
  ambito_id: 1
  children:
    - table: resultados
      cond: "where resultado = :resultado and objetivo_id = :objetivo_id"
      data:
        - resultado: O7R1
          desc_resultado: Unidades educativas com órgãos de gestão plenamente funcionais.
          objetivo_id: objetivo_id()
        - resultado: O7R2
          desc_resultado: Generalização dos Projetos Educativos como instrumento de gestão.
          objetivo_id: objetivo_id()
        - resultado: O7R3
          desc_resultado: Supervisão e acompanhamento pedagógico reforçados.
          objetivo_id: objetivo_id()
        - resultado: O7R4
          desc_resultado: Utilização efetiva do SIGE para apoio à gestão escolar.
          objetivo_id: objetivo_id()
        - resultado: O7R5
          desc_resultado: Melhoria da eficiência e eficácia da gestão escolar.
          objetivo_id: objetivo_id()
    - table: indicadores
      cond: "where indicador = :indicador and objetivo_id = :objetivo_id"
      data:
        - indicador: O7IND35
          desc_indicador: Nº de Unidades Educativas com órgãos de gestão plenamente funcionais
          valor_baseline: 83
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND352026
              ano_id: 1
              ano: 2026
              meta_valor: 83
              indicador_id: indicador_id()
        - indicador: O7IND36
          desc_indicador: Nº de dirigentes das Unidades Educativas capacitados em técnicas de elaboração e gestão de projeto educativo como instrumento estratégico de intervenção
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND362026
              ano_id: 1
              ano: 2026
              meta_valor: 617
              indicador_id: indicador_id()
        - indicador: O7IND37
          desc_indicador: Nº de Unidades Educativas que adotam o projeto educativo como instrumento estratégico de gestão
          valor_baseline: 83
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND372026
              ano_id: 1
              ano: 2026
              meta_valor: 83
              indicador_id: indicador_id()
        - indicador: O7IND38
          desc_indicador: Nº de Unidades Educativas abrangidas com pelo menos uma resultado de supervisão anual
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND382026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
        - indicador: O7IND39
          desc_indicador: Nº de dirigentes e pessoal administrativo das Unidades Educativas capacitados em técnicas de uso e manuseamento do SIGE
          valor_baseline: null
          unidades_id: 2
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND392026
              ano_id: 1
              ano: 2026
              meta_valor: 917
              indicador_id: indicador_id()
        - indicador: O7IND40
          desc_indicador: Percentagem de escolas com SIGE funcional e que reportam dados completos e de qualidade
          valor_baseline: null
          unidades_id: 1
          objetivo_id: objetivo_id()
          children:
            table: metas
            cond: "where ano_id = :ano_id and indicador_id = :indicador_id and objetivo_id = :objetivo_id"
            data:
              meta: O7IND402026
              ano_id: 1
              ano: 2026
              meta_valor: 1.00
              indicador_id: indicador_id()
```

