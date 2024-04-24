export const datas = [
  {
    id: 1,
    name: "User 01",
    role: "admin",
    email: "user01@gmail.com",
    actions: {
      objects: [
        {
          name: "Delete",
          href: "/companies/delete/1454523",
          iconOnly: true,
        },
        {
          name: "Edit",
          href: "/companies/update/1454523",
          iconOnly: true,
        },
      ],
    },
  },
  {
    id: 2,
    name: "User 02",
    role: "admin",
    email: "user02@gmail.com",
    actions: {
      objects: [
        {
          name: "Delete",
          href: "/companies/delete/1454523",
          iconOnly: true,
        },
        {
          name: "Edit",
          href: "/companies/update/1454523",
          iconOnly: true,
        },
      ],
    },
  },
];

export const informationFormDatas = {
  grid: 4,
  title: "General Information",
  titleGrid: 1,
  innerGrid: 4,
  data: [
    {
      name: "User ID",
      value: 1,
      grid: 4,
    },
    {
      name: "Username",
      value: "User 01",
      grid: 2,
    },
    {
      name: "Full Name",
      value: "Admin",
      grid: 2,
    },
    {
      name: "Email",
      value: "user01@gmail.com",
      grid: 2,
    },
    {
      name: "Role",
      value: "Admin",
      grid: 2,
    },
  ],
};

export const configs = {
  headers: "Admin List",
  actions: [
    {
      name: "Add",
      icon: "plus",
      href: "/",
      type: "brand",
    },
  ],
};
