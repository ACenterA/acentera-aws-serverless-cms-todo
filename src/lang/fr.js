export default {
  table: {
    date: 'Date',
    create: 'Creation',
    modification: 'Modification',
    add: 'Ajouter',
    status: 'Status',
    title: 'Titre',
    active: 'Actif',
    select: 'Selectionner',
    action: 'Action',
    delete: 'Effacer',
    cancel: 'Annuler',
    save: 'Sauvegarder',
    confirm: 'Sauvegarder',
    update: 'Mettre a jour',
    edit: 'Modifier'
  },
  plugins: {
    'serverless-cms': {
      button: {
        launch_create: 'Lancer la création',
        launch_create_cluster: 'Lancer la création',
        clear: 'Effacer le formulaire',
        save_state: 'Enregistrer',
        refresh_state: 'Rafraîchir'
      },
      error: {
        not_started: 'La création d\'image pour {title} n\'a pas été démarré.',
        not_started_clusetr: 'La création du cluster {title} n\'a pas été démarré.'
      },
      clustertitle: 'Création du cluster',
      unique_stack_name: 'Le nom doit être unique.',
      stack_state: 'Le status de {name} est {status}.',
      environment: 'Valeur du tag d\'environment ainsi que le préfixe pour les ressources créé',
      'createserverless-cmstitle': 'Création du cluster ECS',
      ecs: 'Informations',
      createecsnodegroup: 'Création de serveurs',
      createecsdiscoverygroup: 'Création du service de discovery',
      task_placement_tags: 'Si défini, les serveurs vont se limiter qu\'aux tâches ayant cette valeur'
    },
    test: {
      amiCreateTitle: 'Création d\'une image personalisé',
      personal_details: 'Détails personnel',
      region_details: 'AMI & Region',
      ami: 'Image ( AMI )'
    }
  },
  error: {
    no_vpc_in_region: 'Aucun VPC\'s disponible dans la region {region}',
    no_subnet_in_region_vpc: 'Aucun résau disponible en {region} pour le vpc {vpc}'
  }
}
