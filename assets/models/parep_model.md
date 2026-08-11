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

## ESTADO_PROGRAMA
```yaml
table: estado_programa
comment: Estado do programa
tooltip: Lista de estados possíveis para um programa.
columns:
  estado_programa_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do programa", tooltip: "Chave primária do estado." }
  estado_programa: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_programa_desc: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_programa_id, order: ASC}]}
```

## ESTADO_OBJETIVO
```yaml
table: estado_objetivo
comment: Estado do objetivo
tooltip: Lista de estados possíveis para um objetivo.
columns:
  estado_objetivo_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do objetivo", tooltip: "Chave primária do estado." }
  estado_objetivo: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_objetivo_desc: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_objetivo_id, order: ASC}]}
```

## ESTADO_FONTE
```yaml
table: estado_fonte
comment: Estado de validação da fonte
tooltip: Estados de validação de um registo face à fonte oficial.
columns:
  estado_fonte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da fonte", tooltip: "Chave primária do estado." }
  estado_fonte: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_fonte_desc: { type: text, comment: "Descrição", tooltip: "Descrição breve do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_fonte_id, order: ASC}]}
data:
  - {estado_fonte: Pendente, estado_fonte_desc: "Ainda não validado"}
  - {estado_fonte: Confirmado, estado_fonte_desc: "Confirmado pela fonte oficial"}
```

## TIPOS_INDICADOR
```yaml
table: tipos_indicador
comment: Tipo de indicador
tooltip: Tipos de valores usados pelos indicadores.
columns:
  tipo_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de indicador", tooltip: "Chave primária do tipo." }
  tipo_indicador: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  tipo_indicador_desc: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de indicador." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o tipo está excluído." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: tipo_indicador_id, order: ASC}]}
data:
  - {tipo_indicador_id: 1, tipo_indicador: Percentagem, tipo_indicador_desc: "Valor percentual"}
  - {tipo_indicador_id: 2, tipo_indicador: Número, tipo_indicador_desc: "Valor quantitativo absoluto"}
  - {tipo_indicador_id: 3, tipo_indicador: Taxa, tipo_indicador_desc: "Taxa ou proporção"}
  - {tipo_indicador_id: 4, tipo_indicador: Rácio, tipo_indicador_desc: "Relação entre valores"}
  - {tipo_indicador_id: 5, tipo_indicador: Valor monetário, tipo_indicador_desc: "Valor monetário"}
  - {tipo_indicador_id: 6, tipo_indicador: Índice, tipo_indicador_desc: "Valor indexado"}
  - {tipo_indicador_id: 7, tipo_indicador: Pontuação, tipo_indicador_desc: "Valor numa escala"}
  - {tipo_indicador_id: 8, tipo_indicador: Sim/Não, tipo_indicador_desc: "Resultado binário"}
  - {tipo_indicador_id: 9, tipo_indicador: Texto, tipo_indicador_desc: "Resultado qualitativo"}
```

## DIMENSOES
```yaml
table: dimensoes
comment: Dimensão de desagregação
tooltip: Dimensões usadas para desagregar as medições.
columns:
  dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da dimensão", tooltip: "Chave primária da dimensão." }
  dimensao: { type: varchar, len: 100, nullable: false, comment: "Nome da dimensão", tooltip: "Nome apresentado ao utilizador." }
  dimensao_desc: { type: text, comment: "Descrição", tooltip: "Descrição da dimensão." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se a dimensão pode ser utilizada." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a dimensão está excluída." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: dimensao_id, order: ASC}]}
```

## VALORES_DIMENSAO
```yaml
table: valores_dimensao
comment: Valor de dimensão
tooltip: Valores possíveis de cada dimensão de desagregação.
columns:
  valor_dimensao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do valor da dimensão", tooltip: "Chave primária do valor." }
  dimensao_id: { type: integer, fk: "dimensoes.dimensao_id", nullable: false, comment: "Dimensão", tooltip: "Dimensão a que o valor pertence.", form_display: true, table_display: true }
  valor_dimensao: { type: varchar, len: 255, nullable: false, comment: "Nome do valor", tooltip: "Nome apresentado ao utilizador.", form_display: true, table_display: true }
  valor_dimensao_ordem: { type: integer, default: 0, comment: "Ordem", tooltip: "Ordem de apresentação." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se o valor está ativo." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o valor está excluído." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: valor_dimensao_ordem, order: ASC}]}
```

## TIPOS_ACAO
```yaml
table: tipos_acao
comment: Tipo de ação
tooltip: Tipos de ações executáveis no programa.
columns:
  tipo_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de ação", tooltip: "Chave primária do tipo." }
  tipo_acao: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  tipo_acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de ação." }
  ativa: { type: boolean, default: true, comment: "Ativa", tooltip: "Indica se o tipo está ativo." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o tipo está excluído." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: tipo_acao_id, order: ASC}]}
data:
  - {tipo_acao: Oficina, tipo_acao_desc: "Oficina ou sessão de trabalho"}
  - {tipo_acao: Formação, tipo_acao_desc: "Formação ou capacitação"}
  - {tipo_acao: Subvenção, tipo_acao_desc: "Transferência financeira"}
  - {tipo_acao: Equipamento, tipo_acao_desc: "Fornecimento de equipamento"}
  - {tipo_acao: Reabilitação, tipo_acao_desc: "Melhoria de espaços"}
  - {tipo_acao: Consultoria, tipo_acao_desc: "Serviço especializado"}
  - {tipo_acao: Estudo, tipo_acao_desc: "Estudo ou avaliação"}
  - {tipo_acao: Aquisição, tipo_acao_desc: "Aquisição de bens ou serviços"}
  - {tipo_acao: Outro, tipo_acao_desc: "Outro tipo de ação"}
```

## PROGRAMAS
```yaml
table: programas
comment: Programa
tooltip: Programa ou instrumento estratégico gerido pelo sistema.
columns:
  programa_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do programa", tooltip: "Chave primária do programa.", form_display: true, table_display: true, order: 1 }
  programa: { type: varchar, len: 255, nullable: false, comment: "Nome", tooltip: "Nome oficial do programa.", form_display: true, table_display: true, form_size: 8, order: 3 }
  programa_desc: { type: text, comment: "Descrição", tooltip: "Descrição geral do programa.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início do programa.", form_display: true, table_display: true, form_size: 4, order: 5 }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data prevista de conclusão.", form_display: true, table_display: true, form_size: 4, order: 6 }
  estado_programa_id: { type: integer, fk: "estado_programa.estado_programa_id", comment: "Estado", tooltip: "Estado atual do programa.", form_display: true, table_display: true, form_size: 4, order: 7 }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se o programa está excluído." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {objetivos: true}
table_layout:
  default_order: [{field: programa_id, order: DESC}]
data:
  - {programa_id: 1, programa: "Partnership Compact Cabo Verde", programa_desc: "Programa nacional para o reforço do capital humano, da empregabilidade e do crescimento inclusivo.", data_inicio: "2024-01-01", data_fim: "2029-12-31", excluded: false}
  - {programa_id: 2, programa: "Transformação Digital da Administração Pública", programa_desc: "Programa de modernização dos serviços públicos e melhoria da qualidade dos dados para decisão.", data_inicio: "2025-01-01", data_fim: "2028-12-31", excluded: false}
```

## OBJETIVOS
```yaml
table: objetivos
comment: Objetivo
tooltip: Objetivo estratégico ou subobjetivo do programa.
columns:
  objetivo_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do objetivo", tooltip: "Chave primária do objetivo.", form_display: true, table_display: true, order: 1 }
  programa_id: { type: integer, fk: "programas.programa_id", nullable: false, comment: "Programa", tooltip: "Programa a que o objetivo pertence.", form_display: true, table_display: true, order: 2 }
  objetivo_pai_id: { type: integer, fk: "objetivos.objetivo_id", comment: "Objetivo pai", tooltip: "Objetivo hierarquicamente superior.", form_display: true, table_display: true }
  objetivo: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do objetivo.", form_display: true, table_display: true, form_size: 8, order: 3 }
  objetivo_desc: { type: text, comment: "Descrição", tooltip: "Descrição do objetivo.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  objetivo_ordem: { type: integer, default: 0, comment: "Ordem", tooltip: "Ordem de apresentação.", form_display: true, table_display: true, form_size: 3, order: 5 }
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
  sub_form_size: 12
  allow_in_subform: {indicadores: true, atividades: true}
table_layout:
  default_order: [{field: objetivo_id, order: DESC}]
data:
  - {objetivo_id: 1, programa_id: 1, objetivo: "Melhorar a qualidade e a relevância da formação profissional", objetivo_desc: "Alinhar a oferta de formação com as necessidades do mercado de trabalho e dos setores prioritários.", objetivo_ordem: 1, excluded: false}
  - {objetivo_id: 2, programa_id: 1, objetivo: "Aumentar a empregabilidade de jovens e mulheres", objetivo_desc: "Promover competências, transição para o emprego e participação económica inclusiva.", objetivo_ordem: 2, excluded: false}
  - {objetivo_id: 3, programa_id: 1, objetivo: "Reforçar a produtividade e a sustentabilidade das empresas", objetivo_desc: "Apoiar empresas na adoção de melhores práticas, tecnologia e modelos de crescimento sustentável.", objetivo_ordem: 3, excluded: false}
  - {objetivo_id: 4, programa_id: 2, objetivo: "Digitalizar serviços públicos prioritários", objetivo_desc: "Disponibilizar serviços públicos mais simples, acessíveis e orientados por dados.", objetivo_ordem: 1, excluded: false}
```

## INDICADORES
```yaml
table: indicadores
comment: Indicador
tooltip: Indicador utilizado para medir um objetivo.
columns:
  indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do indicador", tooltip: "Chave primária do indicador.", form_display: true, table_display: true, order: 1 }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo monitorizado.", form_display: true, table_display: true, order: 2 }
  tipo_indicador_id: { type: integer, fk: "tipos_indicador.tipo_indicador_id", comment: "Tipo", tooltip: "Tipo de valor do indicador.", form_display: true, table_display: true, order: 3 }
  indicador: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do indicador.", form_display: true, table_display: true, form_size: 8, order: 5 }
  indicador_desc: { type: text, comment: "Descrição", tooltip: "Descrição do indicador.", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 6 }
  unidade: { type: varchar, len: 100, comment: "Unidade", tooltip: "Unidade de medida.", form_display: true, table_display: true, form_size: 4, order: 7 }
  valor_base: { type: decimal, len: 20, scale: 6, comment: "Valor de base", tooltip: "Valor numérico de referência.", form_display: true, table_display: true, form_size: 4, order: 8 }
  ano_base: { type: integer, comment: "Ano de base", tooltip: "Ano do valor de referência.", form_display: true, table_display: true, form_size: 4, order: 9 }
  base_texto: { type: text, comment: "Base textual", tooltip: "Valor textual de referência.", form_display: true, table_display: true, form_size: 12, order: 10 }
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
  sub_form_size: 12
  allow_in_subform: {metas_indicadores: true, medicoes_indicadores: true, indicador_dimensoes: true}
  tabs_steps_conf:
    - {label: Geral, fields: [objetivo_id, tipo_indicador_id, indicador, indicador_desc, unidade, direcao_indicador_id]}
    - {label: Base, fields: [valor_base, ano_base, base_texto, metodo_calculo, fonte_dados, frequencia]}
    - {label: Fonte, fields: [estado_fonte_id, referencia_fonte]}
table_layout:
  default_order: [{field: indicador_id, order: DESC}]
data:
  - {indicador_id: 1, objetivo_id: 1, tipo_indicador_id: 2, indicador: "Participantes que concluem formação profissional", indicador_desc: "Número anual de participantes que concluem com aproveitamento uma formação profissional apoiada.", unidade: "pessoas", valor_base: 1800, ano_base: 2023, direcao_indicador_id: 1, frequencia: "Anual", excluded: false}
  - {indicador_id: 2, objetivo_id: 1, tipo_indicador_id: 1, indicador: "Taxa de satisfação dos empregadores com diplomados", indicador_desc: "Percentagem de empregadores que avaliam positivamente as competências dos diplomados.", unidade: "%", valor_base: 62, ano_base: 2023, direcao_indicador_id: 1, frequencia: "Anual", excluded: false}
  - {indicador_id: 3, objetivo_id: 2, tipo_indicador_id: 1, indicador: "Taxa de colocação profissional de jovens formados", indicador_desc: "Percentagem de jovens formados colocados no mercado de trabalho até seis meses após a conclusão.", unidade: "%", valor_base: 38, ano_base: 2023, direcao_indicador_id: 1, frequencia: "Semestral", excluded: false}
  - {indicador_id: 4, objetivo_id: 2, tipo_indicador_id: 1, indicador: "Participação de mulheres em programas de empregabilidade", indicador_desc: "Percentagem de participantes mulheres nos programas de empregabilidade apoiados.", unidade: "%", valor_base: 46, ano_base: 2023, direcao_indicador_id: 1, frequencia: "Anual", excluded: false}
  - {indicador_id: 5, objetivo_id: 3, tipo_indicador_id: 2, indicador: "Empresas apoiadas na adoção de soluções digitais", indicador_desc: "Número de empresas que implementam pelo menos uma solução digital com apoio do programa.", unidade: "empresas", valor_base: 75, ano_base: 2023, direcao_indicador_id: 1, frequencia: "Trimestral", excluded: false}
  - {indicador_id: 6, objetivo_id: 4, tipo_indicador_id: 2, indicador: "Serviços públicos prioritários digitalizados", indicador_desc: "Número de serviços públicos prioritários disponibilizados integralmente em formato digital.", unidade: "serviços", valor_base: 12, ano_base: 2024, direcao_indicador_id: 1, frequencia: "Trimestral", excluded: false}
```

## DIRECAO_INDICADOR
```yaml
table: direcao_indicador
comment: Direção do indicador
tooltip: Direção desejável para interpretar a evolução do indicador.
columns:
  direcao_indicador_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da direção", tooltip: "Chave primária da direção." }
  direcao_indicador: { type: varchar, len: 100, nullable: false, comment: "Nome da direção", tooltip: "Nome apresentado ao utilizador." }
  direcao_indicador_desc: { type: text, comment: "Descrição", tooltip: "Descrição da direção." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {direcao_indicador_id: 1, direcao_indicador: Aumentar, direcao_indicador_desc: "Valores maiores são desejáveis"}
  - {direcao_indicador_id: 2, direcao_indicador: Diminuir, direcao_indicador_desc: "Valores menores são desejáveis"}
  - {direcao_indicador_id: 3, direcao_indicador: Manter, direcao_indicador_desc: "Manter o valor de referência"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: direcao_indicador_id, order: ASC}]}
```

## METAS_INDICADORES
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
  sub_form_size: 12
table_layout:
  default_order: [{field: ano, order: ASC}]
data:
  - {meta_indicador_id: 1, indicador_id: 1, ano: 2025, valor_meta: 2200, notas: "Meta intermédia do programa.", excluded: false}
  - {meta_indicador_id: 2, indicador_id: 1, ano: 2026, valor_meta: 2600, notas: "Meta de expansão da cobertura.", excluded: false}
  - {meta_indicador_id: 3, indicador_id: 2, ano: 2025, valor_meta: 68, meta_texto: "Percentagem", excluded: false}
  - {meta_indicador_id: 4, indicador_id: 3, ano: 2025, valor_meta: 45, meta_texto: "Percentagem", excluded: false}
  - {meta_indicador_id: 5, indicador_id: 4, ano: 2025, valor_meta: 50, meta_texto: "Percentagem", excluded: false}
  - {meta_indicador_id: 6, indicador_id: 5, ano: 2025, valor_meta: 120, excluded: false}
  - {meta_indicador_id: 7, indicador_id: 6, ano: 2025, valor_meta: 20, excluded: false}
```

## MEDICOES_INDICADORES
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
  sub_form_size: 12
table_layout:
  default_order: [{field: ano, order: DESC}]
data:
  - {medicao_indicador_id: 1, indicador_id: 1, ano: 2025, periodo: "S1", valor: 980, fonte: "Relatório de execução do programa", excluded: false}
  - {medicao_indicador_id: 2, indicador_id: 2, ano: 2025, periodo: "S1", valor: 65, fonte: "Inquérito a empregadores", excluded: false}
  - {medicao_indicador_id: 3, indicador_id: 3, ano: 2025, periodo: "S1", valor: 42, fonte: "Sistema de seguimento de diplomados", excluded: false}
  - {medicao_indicador_id: 4, indicador_id: 4, ano: 2025, periodo: "S1", valor: 49, fonte: "Relatório de participação", excluded: false}
  - {medicao_indicador_id: 5, indicador_id: 5, ano: 2025, periodo: "T1", valor: 31, fonte: "Relatório de apoio empresarial", excluded: false}
  - {medicao_indicador_id: 6, indicador_id: 6, ano: 2025, periodo: "T1", valor: 15, fonte: "Catálogo de serviços públicos", excluded: false}
```

## INDICADOR_DIMENSOES
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
table_layout:
  default_order: [{field: indicador_dimensao_id, order: DESC}]
```

## MEDICAO_DIMENSOES
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
table_layout:
  default_order: [{field: medicao_dimensao_id, order: DESC}]
```

## ATIVIDADES
```yaml
table: atividades
comment: Atividade
tooltip: Atividade prevista para realizar um objetivo.
columns:
  atividade_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da atividade", tooltip: "Chave primária da atividade.", form_display: true, table_display: true }
  objetivo_id: { type: integer, fk: "objetivos.objetivo_id", nullable: false, comment: "Objetivo", tooltip: "Objetivo da atividade.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código da atividade.", form_display: true, table_display: true }
  atividade: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da atividade.", form_display: true, table_display: true }
  atividade_desc: { type: text, comment: "Descrição", tooltip: "Descrição da atividade.", form_display: true, table_display: true, form_long_text: true }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início da atividade." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão da atividade." }
  estado_atividade_id: { type: integer, fk: "estado_atividade.estado_atividade_id", comment: "Estado", tooltip: "Estado atual da atividade.", form_display: true, table_display: true }
  responsavel_id: { type: integer, fk: "partes.parte_id", comment: "Responsável", tooltip: "Parte responsável pela atividade." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a atividade está excluída." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {acoes: true}
table_layout:
  default_order: [{field: atividade_id, order: DESC}]
```

## ACOES
```yaml
table: acoes
comment: Ação
tooltip: Ação concreta executada no âmbito de uma atividade.
columns:
  acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da ação", tooltip: "Chave primária da ação.", form_display: true, table_display: true }
  atividade_id: { type: integer, fk: "atividades.atividade_id", nullable: false, comment: "Atividade", tooltip: "Atividade a que a ação pertence.", form_display: true, table_display: true }
  tipo_acao_id: { type: integer, fk: "tipos_acao.tipo_acao_id", comment: "Tipo", tooltip: "Tipo da ação.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 50, nullable: false, comment: "Código", tooltip: "Código da ação.", form_display: true, table_display: true }
  acao: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da ação.", form_display: true, table_display: true }
  acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição da ação.", form_display: true, table_display: true, form_long_text: true }
  data_inicio: { type: date, comment: "Data de início", tooltip: "Data de início da ação." }
  data_fim: { type: date, comment: "Data de fim", tooltip: "Data de conclusão da ação." }
  estado_acao_id: { type: integer, fk: "estado_acao.estado_acao_id", comment: "Estado", tooltip: "Estado atual da ação.", form_display: true, table_display: true }
  responsavel_id: { type: integer, fk: "partes.parte_id", comment: "Responsável", tooltip: "Parte responsável pela ação." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a ação está excluída." }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {orcamentos_acoes: true, beneficiarios_acao: true, acao_documentos: true}
table_layout:
  default_order: [{field: acao_id, order: DESC}]
```

## PARTES
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
form_layout: {tabs_steps: tabs, form_in_popup: false, size: 9, sub_form_size: 12}
table_layout: {default_order: [{field: parte_id, order: DESC}]}
```

## PESSOAS
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 8}
table_layout: {default_order: [{field: pessoa_id, order: DESC}]}
```

## ORGANIZACOES
```yaml
table: organizacoes
comment: Organização participante
tooltip: Organização envolvida no programa.
columns:
  organizacao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da organização", tooltip: "Chave primária da organização.", form_display: true, table_display: true }
  parte_id: { type: integer, fk: "partes.parte_id", nullable: false, unique: true, comment: "Parte", tooltip: "Entidade comum associada." }
  organizacao: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome oficial da organização.", form_display: true, table_display: true }
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 8}
table_layout: {default_order: [{field: organizacao_id, order: DESC}]}
```

## TIPO_PARTE
```yaml
table: tipo_parte
comment: Tipo de parte
tooltip: Tipos de entidades participantes.
columns:
  tipo_parte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do tipo de parte", tooltip: "Chave primária do tipo." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do tipo", tooltip: "Código técnico do tipo." }
  tipo_parte: { type: varchar, len: 100, nullable: false, comment: "Nome do tipo", tooltip: "Nome apresentado ao utilizador." }
  tipo_parte_desc: { type: text, comment: "Descrição", tooltip: "Descrição do tipo de parte." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PERSON, tipo_parte: Pessoa, tipo_parte_desc: "Pessoa singular"}
  - {codigo: ORGANIZATION, tipo_parte: Organização, tipo_parte_desc: "Pessoa coletiva"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: tipo_parte_id, order: ASC}]}
```

## ORCAMENTOS_ACOES
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: ano, order: DESC}]}
```

## BENEFICIARIOS_ACAO
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: beneficiario_acao_id, order: DESC}]}
```

## DOCUMENTOS
```yaml
table: documentos
comment: Documento
tooltip: Metadados de um documento armazenado.
columns:
  documento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do documento", tooltip: "Chave primária do documento.", form_display: true, table_display: true }
  documento: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome descritivo do documento.", form_display: true, table_display: true }
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: documento_id, order: DESC}]}
```

## FONTES
```yaml
table: fontes
comment: Fonte documental
tooltip: Fonte documental utilizada pelo sistema.
columns:
  fonte_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da fonte", tooltip: "Chave primária da fonte.", form_display: true, table_display: true }
  codigo: { type: varchar, len: 100, unique: true, nullable: false, comment: "Código", tooltip: "Código único da fonte.", form_display: true, table_display: true }
  fonte: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome do documento fonte.", form_display: true, table_display: true }
  tipo_documento: { type: varchar, len: 100, comment: "Tipo de documento", tooltip: "Tipo do documento." }
  data_publicacao: { type: date, comment: "Data de publicação", tooltip: "Data de publicação da fonte." }
  url: { type: text, comment: "Endereço eletrónico", tooltip: "Endereço eletrónico da fonte." }
  fonte_desc: { type: text, comment: "Descrição", tooltip: "Descrição da fonte." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a fonte está excluída." }
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: fonte_id, order: DESC}]}
```

## ESTADO_ATIVIDADE
```yaml
# As tabelas seguintes mantêm os estados configuráveis e permitem FKs nos registos operacionais.
table: estado_atividade
comment: Estado da atividade
tooltip: Estados possíveis de uma atividade.
columns:
  estado_atividade_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da atividade", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  estado_atividade: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_atividade_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, estado_atividade: Planeada, estado_atividade_desc: "Atividade planeada"}
  - {codigo: ACTIVE, estado_atividade: Ativa, estado_atividade_desc: "Atividade em execução"}
  - {codigo: COMPLETED, estado_atividade: Concluída, estado_atividade_desc: "Atividade concluída"}
  - {codigo: CANCELLED, estado_atividade: Cancelada, estado_atividade_desc: "Atividade cancelada"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_atividade_id, order: ASC}]}
```

## ESTADO_ACAO
```yaml
table: estado_acao
comment: Estado da ação
tooltip: Estados possíveis de uma ação.
columns:
  estado_acao_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da ação", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  estado_acao: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_acao_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, estado_acao: Planeada, estado_acao_desc: "Ação planeada"}
  - {codigo: ACTIVE, estado_acao: Ativa, estado_acao_desc: "Ação em execução"}
  - {codigo: COMPLETED, estado_acao: Concluída, estado_acao_desc: "Ação concluída"}
  - {codigo: CANCELLED, estado_acao: Cancelada, estado_acao_desc: "Ação cancelada"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_acao_id, order: ASC}]}
```

## CONTRATOS
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
  contrato_desc: { type: text, comment: "Objeto e âmbito", tooltip: "Objeto e âmbito do contrato.", form_display: true, table_display: true, form_long_text: true }
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
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {fases_contrato: true, pagamentos_contrato: true, contrato_documentos: true}
  tabs_steps_conf:
    - {label: Geral, fields: [acao_id, fornecedor_id, numero, titulo, contrato_desc]}
    - {label: Valores, fields: [valor, moeda]}
    - {label: Datas e estado, fields: [data_assinatura, data_inicio, data_fim, estado_contrato_id]}
table_layout:
  default_order: [{field: contrato_id, order: DESC}]
```

## FASES_CONTRATO
```yaml
table: fases_contrato
comment: Fase do contrato
tooltip: Fase de execução de um contrato.
columns:
  fase_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador da fase", tooltip: "Chave primária da fase.", form_display: true, table_display: true }
  contrato_id: { type: integer, fk: "contratos.contrato_id", nullable: false, comment: "Contrato", tooltip: "Contrato a que a fase pertence.", form_display: true, table_display: true }
  fase_contrato: { type: varchar, len: 500, nullable: false, comment: "Nome", tooltip: "Nome da fase.", form_display: true, table_display: true }
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
  sub_form_size: 12
table_layout:
  default_order: [{field: numero_sequencia, order: ASC}]
```

## PAGAMENTOS_CONTRATO
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
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 9
table_layout:
  default_order: [{field: data_pagamento, order: DESC}]
```

## ESTADO_CONTRATO
```yaml
table: estado_contrato
comment: Estado do contrato
tooltip: Estados possíveis de um contrato.
columns:
  estado_contrato_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do contrato", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  estado_contrato: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_contrato_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: DRAFT, estado_contrato: Rascunho, estado_contrato_desc: "Contrato em preparação"}
  - {codigo: ACTIVE, estado_contrato: Ativo, estado_contrato_desc: "Contrato em vigor"}
  - {codigo: COMPLETED, estado_contrato: Concluído, estado_contrato_desc: "Contrato concluído"}
  - {codigo: CANCELLED, estado_contrato: Cancelado, estado_contrato_desc: "Contrato cancelado"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_contrato_id, order: ASC}]}
```

## ESTADO_FASE
```yaml
table: estado_fase
comment: Estado da fase
tooltip: Estados possíveis de uma fase contratual.
columns:
  estado_fase_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado da fase", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  estado_fase: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_fase_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PLANNED, estado_fase: Planeada, estado_fase_desc: "Fase planeada"}
  - {codigo: ACTIVE, estado_fase: Ativa, estado_fase_desc: "Fase em execução"}
  - {codigo: COMPLETED, estado_fase: Concluída, estado_fase_desc: "Fase concluída"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_fase_id, order: ASC}]}
```

## ESTADO_PAGAMENTO
```yaml
table: estado_pagamento
comment: Estado do pagamento
tooltip: Estados possíveis de um pagamento.
columns:
  estado_pagamento_id: { type: integer, pk: true, autoincrement: true, comment: "Identificador do estado do pagamento", tooltip: "Chave primária do estado." }
  codigo: { type: varchar, len: 30, unique: true, nullable: false, comment: "Código do estado", tooltip: "Código técnico do estado." }
  estado_pagamento: { type: varchar, len: 100, nullable: false, comment: "Nome do estado", tooltip: "Nome apresentado ao utilizador." }
  estado_pagamento_desc: { type: text, comment: "Descrição", tooltip: "Descrição do estado." }
  user_id: { type: integer, fk: "users.user_id", comment: "Utilizador", tooltip: "Utilizador responsável pelo registo." }
  created_at: { type: datetime, comment: "Criado em", tooltip: "Data e hora de criação." }
  updated_at: { type: datetime, comment: "Atualizado em", tooltip: "Data e hora da última atualização." }
  excluded: { type: boolean, default: false, comment: "Excluído", tooltip: "Indica se a opção está excluída." }
data:
  - {codigo: PENDING, estado_pagamento: Pendente, estado_pagamento_desc: "Pagamento pendente"}
  - {codigo: PAID, estado_pagamento: Pago, estado_pagamento_desc: "Pagamento efetuado"}
  - {codigo: CANCELLED, estado_pagamento: Cancelado, estado_pagamento_desc: "Pagamento cancelado"}
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 6}
table_layout: {default_order: [{field: estado_pagamento_id, order: ASC}]}
```

## ACAO_DOCUMENTOS
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: acao_documento_id, order: DESC}]}
```

## CONTRATO_DOCUMENTOS
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
form_layout: {tabs_steps: tabs, form_in_popup: true, size: 9}
table_layout: {default_order: [{field: contrato_documento_id, order: DESC}]}
```
