export default {
  table: {
    date: 'Date',
    title: 'Title',
    role: 'Role',
    create: 'Creation',
    modification: 'Modification',
    add: 'Create',
    active: 'Active',
    status: 'Status',
    select: 'Select',
    action: 'Action',
    delete: 'Delete',
    cancel: 'Cancel',
    save: 'Save',
    confirm: 'Save',
    update: 'Update',
    edit: 'Edit'
  },
  post: {
    display_time: 'Date',
    display_time_placeholder: 'Pick date and time',
    title: 'Title',
    userinfo: 'Author',
    author_select: 'Pick an Author',
    author_rating: 'Rating',
    content: 'Description',
    content_placeholder: 'Enter description here',
    description: 'Description',
    description_placeholder: 'Enter description here',
    cover_image: 'Cover Image'
  },
  plugins: {
    'serverless-cms': {
      button: {
        launch_create: 'Create image',
        launch_create_cluster: 'Create cluster',
        clear: 'Clear',
        save_state: 'Save',
        refresh_state: 'Refresh'
      },
      error: {
        not_started: 'Image creation for {title} was not started.',
        not_started_cluster: 'Cluster creation for {title} was not started.'
      },
      clustertitle: 'Create your cluster',
      unique_stack_name: 'The stack name must be unique across the names you already have. (The Environment will be automatically prefixed ie: qa-stackname )',
      stack_state: 'The stack {name} status is {status}.',
      environment: 'Environment tag value and element\s name prefix',
      'createserverless-cmstitle': 'Create an ECS Cluster',
      ecs: 'Informations',
      createecsnodegroup: 'Create new servers',
      createecsdiscoverygroup: 'Create discovery servers',
      task_placement_tags: 'These nodes would only run tasks having that specific tag value, if defined.'
    },
    test: {
      amiCreateTitle: 'Create a customized ami',
      personal_details: 'Personal Details',
      region_details: 'AMI & Region',
      ami: 'Base AMI Selection'
    }
  },
  error: {
    no_vpc_in_region: 'No VPC\'s available in {region}',
    no_subnet_in_region_vpc: 'No subnets available in {region} for vpc {vpc}',
    input: {
      title: 'Missing Display Name',
      email: 'Missing or invalid email address'
    }
  },
  user: {
    title: 'Public Name',
    email: 'E-Mail'
  }
}
