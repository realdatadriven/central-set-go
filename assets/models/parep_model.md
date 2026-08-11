<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
# PAREP_MODEL
```yaml
name: PAREP
description: Modelo de gestão, monitorização e execução do Partnership Compact de Cabo Verde.
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
depends_on: ADMIN
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  PAREP:
    menu_icon: chart-bar
    menu_order: 1
    active: true
    tables:
      - programas
      - objetivos
      - indicadores
      - metas_indicadores
      - medicoes_indicadores
      - atividades
      - acoes
      - orcamentos_acoes
      - partes
      - pessoas
      - organizacoes
      - contratos
      - fases_contrato
      - pagamentos_contrato
      - beneficiarios_acao
      - documentos
      - fontes
  Configuração:
    menu_icon: adjustments
    menu_order: 2
    active: true
    tables:
      - {table: tipos_indicador, active: false}
      - {table: dimensoes, active: false}
      - {table: valores_dimensao, active: false}
      - {table: tipos_acao, active: false}
      - {table: estado_programa, active: false}
      - {table: estado_objetivo, active: false}
      - {table: estado_fonte, active: false}
      - {table: direcao_indicador, active: false}
      - {table: estado_atividade, active: false}
      - {table: estado_acao, active: false}
      - {table: estado_contrato, active: false}
      - {table: estado_fase, active: false}
      - {table: estado_pagamento, active: false}
      - {table: tipo_parte, active: false}
```

## TABELAS DE OPÇÕES

```yaml
table: estado_programa
comment: Estado do programa
tooltip: Lista de estados possíveis para um programa.
columns:
  estado_programa_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do programa", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
```

```yaml
table: estado_objetivo
comment: Estado do objetivo
tooltip: Lista de estados possíveis para um objetivo.
columns:
  estado_objetivo_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do objetivo", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
```

```yaml
table: estado_fonte
comment: Estado de validação da fonte
tooltip: Estados de validação de um registo face à fonte oficial.
columns:
  estado_fonte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da fonte", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: SOURCE_PENDING, nome: Pendente, descricao: "Ainda não validado"}
  - {codigo: CONFIRMED, nome: Confirmado, descricao: "Confirmado pela fonte oficial"}
```

```yaml
table: tipos_indicador
comment: Tipo de indicador
tooltip: Tipos de valores usados pelos indicadores.
columns:
  tipo_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de indicador", tooltip: "Chave primária do tipo." }
  codigo: { type: varchar, len: 50, unique: true, nullable: false, comment: "Código do tipo", tooltip: "Código técnico do tipo." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de indicador." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o tipo está excluído." }
data:
  - {codigo: PERCENTAGE, nome: Percentagem, descricao: "Valor percentual"}
  - {codigo: NUMBER, nome: Número, descricao: "Valor quantitativo absoluto"}
  - {codigo: RATE, nome: Taxa, descricao: "Taxa ou proporção"}
  - {codigo: RATIO, nome: Rácio, descricao: "Relação entre valores"}
  - {codigo: CURRENCY, nome: Valor monetário, descricao: "Valor monetário"}
  - {codigo: INDEX, nome: Índice, descricao: "Valor indexado"}
  - {codigo: SCORE, nome: Pontuação, descricao: "Valor numa escala"}
  - {codigo: BOOLEAN, nome: Sim/Não, descricao: "Resultado binário"}
  - {codigo: TEXT, nome: Texto, descricao: "Resultado qualitativo"}
```

```yaml
table: dimensoes
comment: Dimensão de desagregação
tooltip: Dimensões usadas para desagregar as medições.
columns:
  dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da dimensão", tooltip: "Chave primária da dimensão." }
  codigo: { type: varchar, len: 50, unique: true, nullable: false, comment: "Código da dimensão", tooltip: "Código técnico da dimensão." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome da dimensão", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição da dimensão." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se a dimensão pode ser utilizada." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a dimensão está excluída." }
```

```yaml
table: valores_dimensao
comment: Valor de dimensão
tooltip: Valores possíveis de cada dimensão de desagregação.
columns:
  valor_dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do valor da dimensão", tooltip: "Chave primária do valor." }
  dimensao_id: { type: integer, fk: "dimensoes.dimensao_id", nullable: false, comment: "Dimensão", tooltip: "Dimensão a que o valor pertence.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código do valor", tooltip: "Código técnico do valor.", form_display: true, table_display: true }
  nome: { type: varchar, len: 255, nullable: false, comment: "Nome do valor", tooltip: "Nome apresentado ao utilizador.", form_display: true, table_display: true }
  ordem: { type: integer, default: 0, comment: "Ordem", tooltip: "Ordem de apresentação." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se o valor está ativo." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o valor está excluído." }
```

```yaml
table: tipos_acao
comment: Tipo de ação
tooltip: Tipos de ações executáveis no programa.
columns:
  tipo_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de ação", tooltip: "Chave primária do tipo." }
  codigo: { type: varchar, len: 50, unique: true, nullable: false, comment: "Código do tipo", tooltip: "Código técnico do tipo." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de ação." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se o tipo está ativo." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o tipo está excluído." }
data:
  - {codigo: WORKSHOP, nome: Oficina, descricao: "Oficina ou sessão de trabalho"}
  - {codigo: TRAINING, nome: Formação, descricao: "Formação ou capacitação"}
  - {codigo: GRANT, nome: Subvenção, descricao: "Transferência financeira"}
  - {codigo: EQUIPMENT, nome: Equipamento, descricao: "Fornecimento de equipamento"}
  - {codigo: REHABILITATION, nome: Reabilitação, descricao: "Melhoria de espaços"}
  - {codigo: CONSULTANCY, nome: Consultoria, descricao: "Serviço especializado"}
  - {codigo: STUDY, nome: Estudo, descricao: "Estudo ou avaliação"}
  - {codigo: PROCUREMENT, nome: Aquisição, descricao: "Aquisição de bens ou serviços"}
  - {codigo: OTHER, nome: Outro, descricao: "Outro tipo de ação"}
```

## PROGRAMAS E OBJETIVOS

```yaml
table: programas
comment: Programa
tooltip: Programa ou instrumento estratégico gerido pelo sistema.
columns:
  programa_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do programa", tooltip: "Chave primária do programa.", form_display: true, table_display: true, order: 1 }
  codigo: { type: varchar, len: 50, unique: true, nullable: false, comment: "Código", tooltip: "Código único do programa.", form_display: true, table_display: true, order: 2 }
  nome: { type: varchar, len: 255, nullable: false, comment: "Nome", tooltip: "Nome oficial do programa.", form_display: true, table_display: true, order: 3 }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição geral do programa.", form_display: true, table_display: true, form_long_text: true, order: 4 }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início do programa.", form_display: true, table_display: true }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data prevista de conclusão.", form_display: true, table_display: true }
  estado_programa_id: { type: integer, fk: "estado_programa.estado_programa_id", comment: "Estado", tooltip: "Estado atual do programa.", form_display: true, table_display: true }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o programa está excluído." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: programa_id, order: DESC}]
```

```yaml
table: objetivos
comment: Objetivo
tooltip: Objetivo estratégico ou subobjetivo do programa.
columns:
  objetivo_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do objetivo", tooltip: "Chave primária do objetivo.", form_display: true, table_display: true, order: 1 }
  programa_id: { type: integer, fk: "programas.programa_id", nullable: false, comment: "Programa", tooltip: "Programa a que o objetivo pertence.", form_display: true, table_display: true, order: 2 }
  objetivo_pai_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo pai", tooltip: "Objetivo hierarquicamente superior.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código do objetivo.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do objetivo.", form_display: true, table_display: true }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do objetivo.", form_display: true, table_display: true, form_long_text: true }
  ordem: { type: integer, default: 0, comment: "Ordem", tooltip: "Ordem de apresentação." }
  estado_objetivo_id: { type: integer, fk: "estado_objetivo.estado_objetivo_id", comment: "Estado", tooltip: "Estado atual do objetivo.", form_display: true, table_display: true }
  estado_fonte_id: { type: integer, fk: "estado_fonte.estado_fonte_id", comment: "Validação da fonte", tooltip: "Estado de validação face à fonte oficial.", form_display: true, table_display: true }
  referencia_fonte: { type: text, comment: "Referência da fonte", tooltip: "Referência documental do objetivo." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o objetivo está excluído." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: objetivo_id, order: DESC}]
```

## INDICADORES E MEDIÇÕES

```yaml
table: indicadores
comment: Indicador
tooltip: Indicador utilizado para medir um objetivo.
columns:
  indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do indicador", tooltip: "Chave primária do indicador.", form_display: true, table_display: true, order: 1 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo monitorizado.", form_display: true, table_display: true, order: 2 }
  tipo_indicador_id: { type: integer, fk: "tipos_indicador.tipo_indicador_id", comment: "Tipo", tooltip: "Tipo de valor do indicador.", form_display: true, table_display: true, order: 3 }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código do indicador.", form_display: true, table_display: true, order: 4 }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do indicador.", form_display: true, table_display: true, order: 5 }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do indicador.", form_display: true, table_display: true, form_long_text: true }
  unidade: { type: varchar, len: 100, comment: "Unidade", tooltip: "Unidade de medida.", form_display: true, table_display: true }
  valor_base: { type: decimal, len: 20, scale: 6, comment: "Valor de base", tooltip: "Valor numérico de referência.", form_display: true, table_display: true }
  ano_base: { type: integer, comment: "Ano de base", tooltip: "Ano do valor de referência.", form_display: true, table_display: true }
  base_texto: { type: text, comment: "Base textual", tooltip: "Valor textual de referência.", form_display: true, table_display: true }
  metodo_calculo: { type: text, comment: "Método de cálculo", tooltip: "Metodologia de cálculo." }
  fonte_dados: { type: text, comment: "Fonte dos dados", tooltip: "Origem dos dados." }
  frequencia: { type: varchar, len: 50, comment: "Frequência", tooltip: "Periodicidade de recolha." }
  direcao_indicador_id: { type: integer, fk: "direcao_indicador.direcao_indicador_id", comment: "Direção", tooltip: "Direção desejável do indicador.", form_display: true, table_display: true }
  estado_fonte_id: { type: integer, fk: "estado_fonte.estado_fonte_id", comment: "Validação da fonte", tooltip: "Estado de validação face à fonte oficial." }
  referencia_fonte: { type: text, comment: "Referência da fonte", tooltip: "Referência documental." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o indicador está excluído." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
```

```yaml
table: direcao_indicador
comment: Direção do indicador
tooltip: Direção desejável para interpretar a evolução do indicador.
columns:
  direcao_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da direção", tooltip: "Chave primária da direção." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código da direção", tooltip: "Código técnico da direção." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome da direção", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição da direção." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: INCREASE, nome: Aumentar, descricao: "Valores maiores são desejáveis"}
  - {codigo: DECREASE, nome: Diminuir, descricao: "Valores menores são desejáveis"}
  - {codigo: MAINTAIN, nome: Manter, descricao: "Manter o valor de referência"}
```

```yaml
table: metas_indicadores
comment: Meta anual do indicador
tooltip: Meta definida para um indicador num determinado ano.
columns:
  meta_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da meta", tooltip: "Chave primária da meta.", form_display: true, table_display: true }
  indicador_id: { type: integer, fk: "indicadores.indicador_id", nullable: false, comment: "Indicador", tooltip: "Indicador da meta.", form_display: true, table_display: true }
  ano: { type: integer, nullable: false, comment: "Ano", tooltip: "Ano da meta.", form_display: true, table_display: true }
  valor_meta: { type: decimal, len: 20, scale: 6, comment: "Valor da meta", tooltip: "Meta numérica." }
  meta_texto: { type: text, comment: "Meta textual", tooltip: "Meta qualitativa ou textual." }
  notas: { type: text, comment: "Notas", tooltip: "Observações sobre a meta." }
  referencia_fonte: { type: text, comment: "Referência da fonte", tooltip: "Fonte da meta." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a meta está excluída." }
```

```yaml
table: medicoes_indicadores
comment: Medição do indicador
tooltip: Valor observado de um indicador num período.
columns:
  medicao_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da medição", tooltip: "Chave primária da medição.", form_display: true, table_display: true }
  indicador_id: { type: integer, fk: "indicadores.indicador_id", nullable: false, comment: "Indicador", tooltip: "Indicador medido.", form_display: true, table_display: true }
  ano: { type: integer, nullable: false, comment: "Ano", tooltip: "Ano da medição.", form_display: true, table_display: true }
  periodo: { type: varchar, len: 50, comment: "Período", tooltip: "Período da medição." }
  valor: { type: decimal, len: 20, scale: 6, comment: "Valor", tooltip: "Valor numérico observado." }
  valor_texto: { type: text, comment: "Valor textual", tooltip: "Valor textual observado." }
  fonte: { type: text, comment: "Fonte", tooltip: "Fonte da medição." }
  notas: { type: text, comment: "Notas", tooltip: "Observações da medição." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a medição está excluída." }
```

```yaml
table: indicador_dimensoes
comment: Dimensões do indicador
tooltip: Dimensões aplicáveis à recolha de um indicador.
columns:
  indicador_dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da dimensão do indicador", tooltip: "Chave primária da associação." }
  indicador_id: { type: integer, fk: "indicadores.indicador_id", nullable: false, comment: "Indicador", tooltip: "Indicador associado." }
  dimensao_id: { type: integer, fk: "dimensoes.dimensao_id", nullable: false, comment: "Dimensão", tooltip: "Dimensão associada." }
  obrigatoria: { type: boolean, default: false, comment: "Obrigatória", tooltip: "Indica se a dimensão é obrigatória." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a associação está excluída." }
```

```yaml
table: medicao_dimensoes
comment: Desagregação da medição
tooltip: Valor de dimensão associado a uma medição.
columns:
  medicao_dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da desagregação", tooltip: "Chave primária da desagregação." }
  medicao_indicador_id: { type: integer, fk: "medicoes_indicadores.medicao_indicador_id", nullable: false, comment: "Medição", tooltip: "Medição associada." }
  dimensao_id: { type: integer, fk: "dimensoes.dimensao_id", nullable: false, comment: "Dimensão", tooltip: "Dimensão da desagregação." }
  valor_dimensao_id: { type: integer, fk: "valores_dimensao.valor_dimensao_id", nullable: false, comment: "Valor da dimensão", tooltip: "Valor selecionado para a dimensão." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a desagregação está excluída." }
```

## EXECUÇÃO, PARTES E CONTRATOS

```yaml
table: atividades
comment: Atividade
tooltip: Atividade prevista para realizar um objetivo.
columns:
  atividade_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da atividade", tooltip: "Chave primária da atividade.", form_display: true, table_display: true }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo da atividade.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código da atividade.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da atividade.", form_display: true, table_display: true }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição da atividade.", form_display: true, table_display: true, form_long_text: true }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início da atividade." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão da atividade." }
  estado_atividade_id: { type: integer, fk: "estado_atividade.estado_atividade_id", comment: "Estado", tooltip: "Estado atual da atividade.", form_display: true, table_display: true }
  responsavel_id: { type: integer, fk: "partes.parte_id", comment: "Responsável", tooltip: "Parte responsável pela atividade." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a atividade está excluída." }
```

```yaml
table: acoes
comment: Ação
tooltip: Ação concreta executada no âmbito de uma atividade.
columns:
  acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da ação", tooltip: "Chave primária da ação.", form_display: true, table_display: true }
  atividade_id: { type: integer, fk: "atividades.atividade_id", nullable: false, comment: "Atividade", tooltip: "Atividade a que a ação pertence.", form_display: true, table_display: true }
  tipo_acao_id: { type: integer, fk: "tipos_acao.tipo_acao_id", comment: "Tipo", tooltip: "Tipo da ação.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código da ação.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da ação.", form_display: true, table_display: true }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição da ação.", form_display: true, table_display: true, form_long_text: true }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início da ação." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão da ação." }
  estado_acao_id: { type: integer, fk: "estado_acao.estado_acao_id", comment: "Estado", tooltip: "Estado atual da ação.", form_display: true, table_display: true }
  responsavel_id: { type: integer, fk: "partes.parte_id", comment: "Responsável", tooltip: "Parte responsável pela ação." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a ação está excluída." }
```

```yaml
table: partes
comment: Parte participante
tooltip: Pessoa ou organização participante no programa.
columns:
  parte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da parte", tooltip: "Chave primária da parte.", form_display: true, table_display: true }
  tipo_parte_id: { type: integer, fk: "tipo_parte.tipo_parte_id", nullable: false, comment: "Tipo de parte", tooltip: "Indica se é pessoa ou organização.", form_display: true, table_display: true }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a parte está excluída." }
```

```yaml
table: pessoas
comment: Pessoa participante
tooltip: Pessoa singular participante no programa.
columns:
  pessoa_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da pessoa", tooltip: "Chave primária da pessoa.", form_display: true, table_display: true }
  parte_id: { type: integer, fk: "partes.parte_id", nullable: false, unique: true, comment: "Parte", tooltip: "Entidade comum associada." }
  nome_proprio: { type: varchar, len: 150, nullable: false, comment: "Nome próprio", tooltip: "Nome próprio da pessoa.", form_display: true, table_display: true }
  apelido: { type: varchar, len: 150, nullable: false, comment: "Apelido", tooltip: "Apelido da pessoa.", form_display: true, table_display: true }
  data_nascimento: { type: date, comment: "Data de nascimento", tooltip: "Data de nascimento." }
  sexo: { type: varchar, len: 20, comment: "Sexo", tooltip: "Sexo da pessoa." }
  identificacao_nacional: { type: varchar, len: 100, comment: "Identificação nacional", tooltip: "Número de identificação oficial." }
  telefone: { type: varchar, len: 50, comment: "Telefone", tooltip: "Número de telefone." }
  email: { type: varchar, len: 255, comment: "Correio eletrónico", tooltip: "Endereço de correio eletrónico." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a pessoa está excluída." }
```

```yaml
table: organizacoes
comment: Organização participante
tooltip: Organização envolvida no programa.
columns:
  organizacao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da organização", tooltip: "Chave primária da organização.", form_display: true, table_display: true }
  parte_id: { type: integer, fk: "partes.parte_id", nullable: false, unique: true, comment: "Parte", tooltip: "Entidade comum associada." }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome oficial da organização.", form_display: true, table_display: true }
  tipo_organizacao: { type: varchar, len: 50, comment: "Tipo de organização", tooltip: "Tipo de organização." }
  numero_fiscal: { type: varchar, len: 100, comment: "Número fiscal", tooltip: "Número fiscal da organização." }
  endereco: { type: text, comment: "Endereço", tooltip: "Endereço da organização." }
  telefone: { type: varchar, len: 50, comment: "Telefone", tooltip: "Número de telefone." }
  email: { type: varchar, len: 255, comment: "Correio eletrónico", tooltip: "Endereço de correio eletrónico." }
  sitio_web: { type: varchar, len: 500, comment: "Sítio Web", tooltip: "Endereço do sítio Web." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a organização está excluída." }
```

```yaml
table: tipo_parte
comment: Tipo de parte
tooltip: Tipos de entidades participantes.
columns:
  tipo_parte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de parte", tooltip: "Chave primária do tipo." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do tipo", tooltip: "Código técnico do tipo." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de parte." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PERSON, nome: Pessoa, descricao: "Pessoa singular"}
  - {codigo: ORGANIZATION, nome: Organização, descricao: "Pessoa coletiva"}
```

```yaml
table: orcamentos_acoes
comment: Orçamento anual da ação
tooltip: Valores orçamentais planeados e executados por ano.
columns:
  orcamento_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do orçamento", tooltip: "Chave primária do orçamento.", form_display: true, table_display: true }
  acao_id: { type: integer, fk: "acoes.acao_id", nullable: false, comment: "Ação", tooltip: "Ação orçamentada.", form_display: true, table_display: true }
  ano: { type: integer, nullable: false, comment: "Ano", tooltip: "Ano orçamental.", form_display: true, table_display: true }
  moeda: { type: varchar, len: 3, default: CVE, nullable: false, comment: "Moeda", tooltip: "Código da moeda." }
  valor_estimado: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor estimado", tooltip: "Valor estimado." }
  valor_aprovado: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor aprovado", tooltip: "Valor aprovado." }
  valor_comprometido: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor comprometido", tooltip: "Valor comprometido." }
  valor_executado: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor executado", tooltip: "Valor executado." }
  notas: { type: text, comment: "Notas", tooltip: "Observações do orçamento." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o orçamento está excluído." }
```

```yaml
table: beneficiarios_acao
comment: Beneficiário da ação
tooltip: Parte beneficiária de uma ação.
columns:
  beneficiario_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do beneficiário", tooltip: "Chave primária do beneficiário.", form_display: true, table_display: true }
  acao_id: { type: integer, fk: "acoes.acao_id", nullable: false, comment: "Ação", tooltip: "Ação que beneficia a parte." }
  parte_id: { type: integer, fk: "partes.parte_id", nullable: false, comment: "Parte", tooltip: "Parte beneficiária." }
  papel_beneficiario: { type: varchar, len: 100, comment: "Papel do beneficiário", tooltip: "Papel desempenhado pelo beneficiário." }
  quantidade_estimada: { type: integer, comment: "Quantidade estimada", tooltip: "Número estimado de beneficiários." }
  quantidade_real: { type: integer, comment: "Quantidade real", tooltip: "Número efetivamente alcançado." }
  notas: { type: text, comment: "Notas", tooltip: "Observações sobre o beneficiário." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o beneficiário está excluído." }
```

## DOCUMENTOS E FONTES

```yaml
table: documentos
comment: Documento
tooltip: Metadados de um documento armazenado.
columns:
  documento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do documento", tooltip: "Chave primária do documento.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome descritivo do documento.", form_display: true, table_display: true }
  nome_ficheiro: { type: varchar, len: 500, nullable: false, comment: "Nome do ficheiro", tooltip: "Nome original do ficheiro." }
  tipo_mime: { type: varchar, len: 150, comment: "Tipo MIME", tooltip: "Tipo MIME do ficheiro." }
  tamanho: { type: integer, comment: "Tamanho", tooltip: "Tamanho do ficheiro em bytes." }
  caminho_armazenamento: { type: text, nullable: false, comment: "Caminho de armazenamento", tooltip: "Localização do ficheiro." }
  soma_verificacao: { type: varchar, len: 128, comment: "Soma de verificação", tooltip: "Hash de integridade do ficheiro." }
  carregado_em: { type: datetime, comment: "Carregado em", tooltip: "Data e hora de carregamento." }
  carregado_por: { type: integer, fk: "users.user_id", comment: "Carregado por", tooltip: "Utilizador que carregou o documento." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o documento está excluído." }
```

```yaml
table: fontes
comment: Fonte documental
tooltip: Fonte documental utilizada pelo sistema.
columns:
  fonte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da fonte", tooltip: "Chave primária da fonte.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 100, unique: true, nullable: false, comment: "Código", tooltip: "Código único da fonte.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do documento fonte.", form_display: true, table_display: true }
  tipo_documento: { type: varchar, len: 100, comment: "Tipo de documento", tooltip: "Tipo do documento." }
  data_publicacao: { type: date, comment: "Data de publicação", tooltip: "Data de publicação da fonte." }
  url: { type: text, comment: "Endereço eletrónico", tooltip: "Endereço eletrónico da fonte." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição da fonte." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a fonte está excluída." }
```

## TABELAS DE ESTADOS OPERACIONAIS

```yaml
# As tabelas seguintes mantêm os estados configuráveis e permitem FKs nos registos operacionais.
table: estado_atividade
comment: Estado da atividade
tooltip: Estados possíveis de uma atividade.
columns:
  estado_atividade_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da atividade", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, nome: Planeada, descricao: "Atividade planeada"}
  - {codigo: ACTIVE, nome: Ativa, descricao: "Atividade em execução"}
  - {codigo: COMPLETED, nome: Concluída, descricao: "Atividade concluída"}
  - {codigo: CANCELLED, nome: Cancelada, descricao: "Atividade cancelada"}
```

```yaml
table: estado_acao
comment: Estado da ação
tooltip: Estados possíveis de uma ação.
columns:
  estado_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da ação", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, nome: Planeada, descricao: "Ação planeada"}
  - {codigo: ACTIVE, nome: Ativa, descricao: "Ação em execução"}
  - {codigo: COMPLETED, nome: Concluída, descricao: "Ação concluída"}
  - {codigo: CANCELLED, nome: Cancelada, descricao: "Ação cancelada"}
```

```yaml
table: contratos
comment: Contrato
tooltip: Contrato celebrado para executar uma ação.
columns:
  contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do contrato", tooltip: "Chave primária do contrato.", form_display: true, table_display: true }
  acao_id: { type: integer, fk: "acoes.acao_id", nullable: false, comment: "Ação", tooltip: "Ação executada pelo contrato.", form_display: true, table_display: true }
  fornecedor_id: { type: integer, fk: "partes.parte_id", nullable: false, comment: "Fornecedor", tooltip: "Fornecedor ou prestador do contrato.", form_display: true, table_display: true }
  numero: { type: varchar, len: 100, unique: true, nullable: false, comment: "Número do contrato", tooltip: "Número único do contrato.", form_display: true, table_display: true }
  titulo: { type: varchar, len: 500, nullable: false, comment: "Título", tooltip: "Título do contrato.", form_display: true, table_display: true }
  descricao: { type: text, comment: "Objeto e âmbito", tooltip: "Objeto e âmbito do contrato.", form_display: true, table_display: true, form_long_text: true }
  valor: { type: decimal, len: 20, scale: 2, nullable: false, comment: "Valor", tooltip: "Valor total do contrato.", form_display: true, table_display: true }
  moeda: { type: varchar, len: 3, default: CVE, nullable: false, comment: "Moeda", tooltip: "Código da moeda." }
  data_assinatura: { type: date, comment: "Data de assinatura", tooltip: "Data de assinatura do contrato." }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início do contrato." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão do contrato." }
  estado_contrato_id: { type: integer, fk: "estado_contrato.estado_contrato_id", comment: "Estado", tooltip: "Estado atual do contrato.", form_display: true, table_display: true }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o contrato está excluído." }
```

```yaml
table: fases_contrato
comment: Fase do contrato
tooltip: Fase de execução de um contrato.
columns:
  fase_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da fase", tooltip: "Chave primária da fase.", form_display: true, table_display: true }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato a que a fase pertence.", form_display: true, table_display: true }
  nome: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da fase.", form_display: true, table_display: true }
  numero_sequencia: { type: integer, nullable: false, comment: "Número de sequência", tooltip: "Ordem da fase." }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início da fase." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão da fase." }
  valor_planeado: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor planeado", tooltip: "Valor planeado para a fase." }
  valor_real: { type: decimal, len: 20, scale: 2, default: 0, comment: "Valor real", tooltip: "Valor executado na fase." }
  percentagem_concluida: { type: decimal, len: 5, scale: 2, default: 0, comment: "Percentagem concluída", tooltip: "Percentagem concluída da fase." }
  estado_fase_id: { type: integer, fk: "estado_fase.estado_fase_id", comment: "Estado", tooltip: "Estado atual da fase.", form_display: true, table_display: true }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a fase está excluída." }
```

```yaml
table: pagamentos_contrato
comment: Pagamento do contrato
tooltip: Pagamento associado a um contrato.
columns:
  pagamento_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do pagamento", tooltip: "Chave primária do pagamento.", form_display: true, table_display: true }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato pago.", form_display: true, table_display: true }
  fase_contrato_id: { type: integer, fk: "fases_contrato.fase_contrato_id", comment: "Fase contratual", tooltip: "Fase a que o pagamento se refere." }
  data_pagamento: { type: date, nullable: false, comment: "Data do pagamento", tooltip: "Data do pagamento.", form_display: true, table_display: true }
  valor: { type: decimal, len: 20, scale: 2, nullable: false, comment: "Valor", tooltip: "Valor do pagamento.", form_display: true, table_display: true }
  moeda: { type: varchar, len: 3, default: CVE, nullable: false, comment: "Moeda", tooltip: "Código da moeda." }
  referencia_pagamento: { type: varchar, len: 150, comment: "Referência do pagamento", tooltip: "Referência bancária ou administrativa." }
  estado_pagamento_id: { type: integer, fk: "estado_pagamento.estado_pagamento_id", comment: "Estado", tooltip: "Estado atual do pagamento.", form_display: true, table_display: true }
  notas: { type: text, comment: "Notas", tooltip: "Observações do pagamento." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o pagamento está excluído." }
```

```yaml
table: estado_contrato
comment: Estado do contrato
tooltip: Estados possíveis de um contrato.
columns:
  estado_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do contrato", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: DRAFT, nome: Rascunho, descricao: "Contrato em preparação"}
  - {codigo: ACTIVE, nome: Ativo, descricao: "Contrato em vigor"}
  - {codigo: COMPLETED, nome: Concluído, descricao: "Contrato concluído"}
  - {codigo: CANCELLED, nome: Cancelado, descricao: "Contrato cancelado"}
```

```yaml
table: estado_fase
comment: Estado da fase
tooltip: Estados possíveis de uma fase contratual.
columns:
  estado_fase_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da fase", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, nome: Planeada, descricao: "Fase planeada"}
  - {codigo: ACTIVE, nome: Ativa, descricao: "Fase em execução"}
  - {codigo: COMPLETED, nome: Concluída, descricao: "Fase concluída"}
```

```yaml
table: estado_pagamento
comment: Estado do pagamento
tooltip: Estados possíveis de um pagamento.
columns:
  estado_pagamento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do pagamento", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  nome: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  descricao: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PENDING, nome: Pendente, descricao: "Pagamento pendente"}
  - {codigo: PAID, nome: Pago, descricao: "Pagamento efetuado"}
  - {codigo: CANCELLED, nome: Cancelado, descricao: "Pagamento cancelado"}
```

## ASSOCIAÇÕES DE DOCUMENTOS

```yaml
table: acao_documentos
comment: Documento da ação
tooltip: Associação entre uma ação e os seus documentos.
columns:
  acao_documento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da associação", tooltip: "Chave primária da associação." }
  acao_id: { type: integer, fk: "acoes.acao_id", nullable: false, comment: "Ação", tooltip: "Ação associada." }
  documento_id: { type: integer, fk: "documentos.documento_id", nullable: false, comment: "Documento", tooltip: "Documento associado." }
  tipo_documento: { type: varchar, len: 100, comment: "Tipo de documento", tooltip: "Tipo do documento associado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a associação está excluída." }
```

```yaml
table: contrato_documentos
comment: Documento do contrato
tooltip: Associação entre um contrato e os seus documentos.
columns:
  contrato_documento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da associação", tooltip: "Chave primária da associação." }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato associado." }
  documento_id: { type: integer, fk: "documentos.documento_id", nullable: false, comment: "Documento", tooltip: "Documento associado." }
  tipo_documento: { type: varchar, len: 100, comment: "Tipo de documento", tooltip: "Tipo do documento associado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a associação está excluída." }
```
